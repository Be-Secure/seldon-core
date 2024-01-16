/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed BY
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package io.seldon.dataflow.kafka

import io.seldon.mlops.chainer.ChainerOuterClass.PipelineStepUpdate
import org.apache.kafka.clients.admin.Admin
import org.apache.kafka.clients.admin.NewTopic
import org.apache.kafka.common.KafkaFuture
import org.apache.kafka.common.errors.TopicExistsException
import java.util.concurrent.ExecutionException
import io.klogging.logger as coLogger

class KafkaAdmin(
    adminConfig: KafkaAdminProperties,
    private val streamsConfig: KafkaStreamsParams,
) {
    private val adminClient = Admin.create(adminConfig)

    suspend fun ensureTopicsExist(
        steps: List<PipelineStepUpdate>,
    ) {
        steps
            .flatMap { step -> step.sourcesList + step.sink + step.triggersList }
            .map { topicName -> parseSource(topicName).first }
            .toSet()
            .also {
                logger.info("Topics found are $it")
            }
            .map { topicName ->
                NewTopic(
                    topicName,
                    streamsConfig.numPartitions,
                    streamsConfig.replicationFactor.toShort(),
                ).configs(
                    kafkaTopicConfig(
                        streamsConfig.maxMessageSizeBytes,
                    ),
                )
            }
            .run { adminClient.createTopics(this) }
            .values()
            .also { topicCreations ->
                topicCreations.entries.forEach { creationResult ->
                    awaitKafkaResult(creationResult)
                }
            }
    }

    private suspend fun awaitKafkaResult(result: Map.Entry<String, KafkaFuture<Void>>) {
        try {
            result.value.get()
            logger.info("Topic created ${result.key}")
        } catch (e: ExecutionException) {
            if (e.cause is TopicExistsException) {
                logger.info("Topic already exists ${result.key}")
            } else {
                throw e
            }
        }
    }

    companion object {
        private val logger = coLogger(KafkaAdmin::class)
    }
}
