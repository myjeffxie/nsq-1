package nsqd

import (
	"crypto/md5"
	"crypto/tls"
	"github.com/absolute8511/nsq/internal/levellogger"
	"hash/crc32"
	"io"
	"log"
	"os"
	"time"
)

const (
	MAX_NODE_ID = 1024 * 1024
)

type Options struct {
	// basic options
	ID                         int64         `flag:"worker-id" cfg:"id"`
	Verbose                    bool          `flag:"verbose"`
	ClusterID                  string        `flag:"cluster-id"`
	ClusterLeadershipAddresses string        `flag:"cluster-leadership-addresses" cfg:"cluster_leadership_addresses"`
	TCPAddress                 string        `flag:"tcp-address"`
	RPCPort                    string        `flag:"rpc-port"`
	ReverseProxyPort           string        `flag:"reverse-proxy-port"`
	HTTPAddress                string        `flag:"http-address"`
	HTTPSAddress               string        `flag:"https-address"`
	BroadcastAddress           string        `flag:"broadcast-address"`
	BroadcastInterface         string        `flag:"broadcast-interface"`
	NSQLookupdTCPAddresses     []string      `flag:"lookupd-tcp-address" cfg:"nsqlookupd_tcp_addresses"`
	AuthHTTPAddresses          []string      `flag:"auth-http-address" cfg:"auth_http_addresses"`
	LookupPingInterval         time.Duration `flag:"lookup-ping-interval" arg:"5s"`

	// diskqueue options
	DataPath        string        `flag:"data-path"`
	MemQueueSize    int64         `flag:"mem-queue-size"`
	MaxBytesPerFile int64         `flag:"max-bytes-per-file"`
	SyncEvery       int64         `flag:"sync-every"`
	SyncTimeout     time.Duration `flag:"sync-timeout"`

	QueueScanInterval        time.Duration
	QueueScanRefreshInterval time.Duration
	QueueScanSelectionCount  int
	QueueScanWorkerPoolMax   int
	QueueScanDirtyPercent    float64

	// msg and command options
	MsgTimeout    time.Duration `flag:"msg-timeout" arg:"60s"`
	MaxMsgTimeout time.Duration `flag:"max-msg-timeout"`
	MaxMsgSize    int64         `flag:"max-msg-size" deprecated:"max-message-size" cfg:"max_msg_size"`
	MaxBodySize   int64         `flag:"max-body-size"`
	MaxReqTimeout time.Duration `flag:"max-req-timeout"`
	MaxConfirmWin int64         `flag:"max-confirm-win"`
	ClientTimeout time.Duration

	// client overridable configuration options
	MaxHeartbeatInterval   time.Duration `flag:"max-heartbeat-interval"`
	MaxRdyCount            int64         `flag:"max-rdy-count"`
	MaxOutputBufferSize    int64         `flag:"max-output-buffer-size"`
	MaxOutputBufferTimeout time.Duration `flag:"max-output-buffer-timeout"`

	// statsd integration
	StatsdAddress  string        `flag:"statsd-address"`
	StatsdPrefix   string        `flag:"statsd-prefix"`
	StatsdProtocol string        `flag:"statsd-protocol"`
	StatsdInterval time.Duration `flag:"statsd-interval" arg:"60s"`
	StatsdMemStats bool          `flag:"statsd-mem-stats"`

	// e2e message latency
	E2EProcessingLatencyWindowTime  time.Duration `flag:"e2e-processing-latency-window-time"`
	E2EProcessingLatencyPercentiles []float64     `flag:"e2e-processing-latency-percentile" cfg:"e2e_processing_latency_percentiles"`

	// TLS config
	TLSCert             string `flag:"tls-cert"`
	TLSKey              string `flag:"tls-key"`
	TLSClientAuthPolicy string `flag:"tls-client-auth-policy"`
	TLSRootCAFile       string `flag:"tls-root-ca-file"`
	TLSRequired         int    `flag:"tls-required"`
	TLSMinVersion       uint16 `flag:"tls-min-version"`

	// compression
	DeflateEnabled  bool `flag:"deflate"`
	MaxDeflateLevel int  `flag:"max-deflate-level"`
	SnappyEnabled   bool `flag:"snappy"`

	LogLevel int32  `flag:"log-level" cfg:"log_level"`
	LogDir   string `flag:"log-dir" cfg:"log_dir"`
	Logger   levellogger.Logger
}

func NewOptions() *Options {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	h := md5.New()
	io.WriteString(h, hostname)
	defaultID := int64(crc32.ChecksumIEEE(h.Sum(nil)) % MAX_NODE_ID)

	return &Options{
		ID: defaultID,

		ClusterID:                  "nsq-clusterid-test-only",
		ClusterLeadershipAddresses: "",
		TCPAddress:                 "0.0.0.0:4150",
		HTTPAddress:                "0.0.0.0:4151",
		HTTPSAddress:               "0.0.0.0:4152",
		BroadcastAddress:           hostname,
		BroadcastInterface:         "eth0",

		NSQLookupdTCPAddresses: make([]string, 0),
		AuthHTTPAddresses:      make([]string, 0),
		LookupPingInterval:     5 * time.Second,

		MemQueueSize:    10000,
		MaxBytesPerFile: 100 * 1024 * 1024,
		SyncEvery:       2500,
		SyncTimeout:     2 * time.Second,

		QueueScanInterval:        500 * time.Millisecond,
		QueueScanRefreshInterval: 5 * time.Second,
		QueueScanSelectionCount:  20,
		QueueScanWorkerPoolMax:   4,
		QueueScanDirtyPercent:    0.25,

		MsgTimeout:    60 * time.Second,
		MaxMsgTimeout: 15 * time.Minute,
		MaxMsgSize:    1024 * 1024,
		MaxBodySize:   5 * 1024 * 1024,
		MaxReqTimeout: 1 * time.Hour,
		ClientTimeout: 60 * time.Second,

		MaxHeartbeatInterval:   60 * time.Second,
		MaxRdyCount:            2500,
		MaxOutputBufferSize:    64 * 1024,
		MaxOutputBufferTimeout: 1 * time.Second,
		MaxConfirmWin:          500,

		StatsdPrefix:   "nsq.%s",
		StatsdProtocol: "udp",
		StatsdInterval: 60 * time.Second,
		StatsdMemStats: true,

		E2EProcessingLatencyWindowTime: time.Duration(10 * time.Minute),

		DeflateEnabled:  true,
		MaxDeflateLevel: 6,
		SnappyEnabled:   true,

		TLSMinVersion: tls.VersionTLS10,

		LogLevel: levellogger.LOG_INFO,
		LogDir:   "",
		Logger:   &levellogger.GLogger{},
	}
}
