/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// NetworkMessageSentTotal is total number of sent messages metric
var NetworkMessageSentTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name:      "message_sent_total",
	Help:      "Total number of sent messages",
	Namespace: insolarNamespace,
	Subsystem: "network",
})

// NetworkPacketSentTotal is total number of sent packets metric
var NetworkPacketSentTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name:      "packet_sent_total",
	Help:      "Total number of sent packets",
	Namespace: insolarNamespace,
	Subsystem: "network",
}, []string{"packetType"})

// NetworkPacketReceivedTotal is is total number of received packets metric
var NetworkPacketReceivedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name:      "packet_received_total",
	Help:      "Total number of received packets",
	Namespace: insolarNamespace,
	Subsystem: "network",
}, []string{"packetType"})

// NetworkFutures is current network transport futures count metric
var NetworkFutures = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "futures",
	Help:      "Current network transport futures count",
	Namespace: insolarNamespace,
	Subsystem: "network",
}, []string{"packetType"})
