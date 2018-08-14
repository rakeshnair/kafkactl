package kafkactl

import (
	"testing"

	"sort"

	"github.com/stretchr/testify/assert"
)

//go:generate mockery -dir . -outpkg kafkactl -inpkg -output . -case underscore -name=ClusterAPI

func TestUniformDistStrategy_Assignments(t *testing.T) {
	defer func() { randomStartIndexFn = randomStartIndex }()

	t.Run("singe_topic", func(t *testing.T) {
		topics := []string{"kafka-test-1"}
		expected := []PartitionReplicas{{"kafka-test-1", 0, []BrokerID{1, 3}}, {"kafka-test-1", 1, []BrokerID{3, 2}}, {"kafka-test-1", 2, []BrokerID{2, 1}}}

		mockCluster := &MockClusterAPI{}
		uds := &UniformDistStrategy{
			topics:  topics,
			cluster: mockCluster,
		}
		randomStartIndexFn = func(max int) int { return 0 }
		mockCluster.On("Brokers").Return([]Broker{{Id: 1, Rack: "1"}, {Id: 2, Rack: "1"}, {Id: 3, Rack: "2"}}, nil)
		mockCluster.On("DescribeTopic", "kafka-test-1").Return(
			[]TopicPartitionInfo{
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 0}, Replication: 2, Leader: 1},
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 1}, Replication: 2, Leader: 2},
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 2}, Replication: 2, Leader: 3},
			}, nil)

		actual, err := uds.Assignments()
		if assert.NoError(t, err) {
			sort.Sort(byPartitionInPartitionReplicas(expected))
			assert.EqualValues(t, expected, actual)
		}
	})

	t.Run("two_topics", func(t *testing.T) {
		topics := []string{"kafka-test-1", "kafka-test-2"}
		expected := []PartitionReplicas{
			{"kafka-test-1", 0, []BrokerID{1, 3}},
			{"kafka-test-1", 1, []BrokerID{3, 2}},
			{"kafka-test-1", 2, []BrokerID{2, 1}},
			{"kafka-test-2", 0, []BrokerID{1, 3}},
			{"kafka-test-2", 1, []BrokerID{3, 2}}}

		mockCluster := &MockClusterAPI{}
		uds := &UniformDistStrategy{
			topics:  topics,
			cluster: mockCluster,
		}
		randomStartIndexFn = func(max int) int { return 0 }
		mockCluster.On("Brokers").Return([]Broker{{Id: 1, Rack: "1"}, {Id: 2, Rack: "1"}, {Id: 3, Rack: "2"}}, nil)
		mockCluster.On("DescribeTopic", "kafka-test-1").Return(
			[]TopicPartitionInfo{
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 0}, Replication: 2, Leader: 1},
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 1}, Replication: 2, Leader: 2},
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 2}, Replication: 2, Leader: 3},
			}, nil)
		mockCluster.On("DescribeTopic", "kafka-test-2").Return(
			[]TopicPartitionInfo{
				{TopicPartition: TopicPartition{Topic: "kafka-test-2", Partition: 0}, Replication: 2, Leader: 1},
				{TopicPartition: TopicPartition{Topic: "kafka-test-2", Partition: 1}, Replication: 2, Leader: 2},
			}, nil)

		actual, err := uds.Assignments()
		if assert.NoError(t, err) {
			sort.Sort(byPartitionInPartitionReplicas(expected))
			assert.EqualValues(t, expected, actual)
		}
	})
}

func TestUniformDistStrategy_topicAssignments(t *testing.T) {

	t.Run("no_rack_2_node_cluster", func(t *testing.T) {
		topic := "kafka-test-1"
		expected := []PartitionReplicas{{"kafka-test-1", 0, []BrokerID{1, 2}}, {"kafka-test-1", 1, []BrokerID{2, 1}}}

		mockCluster := &MockClusterAPI{}
		uds := &UniformDistStrategy{
			cluster: mockCluster,
		}
		mockCluster.On("Brokers").Return([]Broker{{Id: 1, Rack: "1"}, {Id: 2, Rack: "1"}}, nil)
		mockCluster.On("DescribeTopic", "kafka-test-1").Return(
			[]TopicPartitionInfo{
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 0}, Replication: 2, Leader: 1},
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 1}, Replication: 2, Leader: 2},
			}, nil)

		actual, err := uds.topicAssignments(topic)
		if assert.NoError(t, err) {
			sort.Sort(byPartitionInPartitionReplicas(expected))
			assert.EqualValues(t, actual, expected)
		}
	})

	t.Run("rack_aware_3_node_cluster", func(t *testing.T) {
		topic := "kafka-test-1"
		expected := []PartitionReplicas{{"kafka-test-1", 0, []BrokerID{1, 3}}, {"kafka-test-1", 1, []BrokerID{3, 2}}, {"kafka-test-1", 2, []BrokerID{2, 1}}}

		mockCluster := &MockClusterAPI{}
		uds := &UniformDistStrategy{
			cluster: mockCluster,
		}
		mockCluster.On("Brokers").Return([]Broker{{Id: 1, Rack: "1"}, {Id: 2, Rack: "1"}, {Id: 3, Rack: "2"}}, nil)
		mockCluster.On("DescribeTopic", "kafka-test-1").Return(
			[]TopicPartitionInfo{
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 0}, Replication: 2, Leader: 1},
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 1}, Replication: 2, Leader: 2},
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 2}, Replication: 2, Leader: 3},
			}, nil)

		actual, err := uds.topicAssignments(topic)
		if assert.NoError(t, err) {
			sort.Sort(byPartitionInPartitionReplicas(expected))
			assert.EqualValues(t, expected, actual)
		}
	})

	t.Run("rack_aware_3_node_cluster_index_1", func(t *testing.T) {
		topic := "kafka-test-1"
		expected := []PartitionReplicas{{"kafka-test-1", 0, []BrokerID{3, 2}}, {"kafka-test-1", 1, []BrokerID{2, 1}}, {"kafka-test-1", 2, []BrokerID{1, 3}}}

		mockCluster := &MockClusterAPI{}
		uds := &UniformDistStrategy{
			cluster: mockCluster,
		}
		mockCluster.On("Brokers").Return([]Broker{{Id: 1, Rack: "1"}, {Id: 2, Rack: "1"}, {Id: 3, Rack: "2"}}, nil)
		mockCluster.On("DescribeTopic", "kafka-test-1").Return(
			[]TopicPartitionInfo{
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 0}, Replication: 2, Leader: 3},
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 1}, Replication: 2, Leader: 2},
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 2}, Replication: 2, Leader: 1},
			}, nil)

		actual, err := uds.topicAssignments(topic)
		if assert.NoError(t, err) {
			sort.Sort(byPartitionInPartitionReplicas(expected))
			assert.EqualValues(t, expected, actual)
		}
	})

	t.Run("rack_aware_3_node_cluster_replication_3", func(t *testing.T) {
		topic := "kafka-test-1"
		expected := []PartitionReplicas{{"kafka-test-1", 0, []BrokerID{1, 3, 2}}, {"kafka-test-1", 1, []BrokerID{3, 2, 1}}, {"kafka-test-1", 2, []BrokerID{2, 1, 3}}}

		mockCluster := &MockClusterAPI{}
		uds := &UniformDistStrategy{
			cluster: mockCluster,
		}
		mockCluster.On("Brokers").Return([]Broker{{Id: 1, Rack: "1"}, {Id: 2, Rack: "1"}, {Id: 3, Rack: "2"}}, nil)
		mockCluster.On("DescribeTopic", "kafka-test-1").Return(
			[]TopicPartitionInfo{
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 0}, Replication: 3, Leader: 1},
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 1}, Replication: 3, Leader: 1},
				{TopicPartition: TopicPartition{Topic: "kafka-test-1", Partition: 2}, Replication: 3, Leader: 1},
			}, nil)

		actual, err := uds.topicAssignments(topic)
		if assert.NoError(t, err) {
			sort.Sort(byPartitionInPartitionReplicas(expected))
			assert.EqualValues(t, expected, actual)
		}
	})
}
