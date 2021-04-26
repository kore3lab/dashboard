package model

import "time"

type Summary struct {
	Node struct {
		Nodename         string            `json:"nodeName"`
		Systemcontainers []ContainerStruct `json:"systemContainers"`
		Starttime        time.Time         `json:"startTime"`
		CPU              CPUStruct         `json:"cpu"`
		Memory           MemoryStruct      `json:"memory"`
		Network          NetworkStruct     `json:"network"`
		Fs               FsStruct          `json:"fs"`
		Runtime          RuntimeStruct     `json:"runtime"`
		Rlimit           RlimitStruct      `json:"rlimit"`
	} `json:"node"`
	Pods []struct {
		Podref           PodrefStruct       `json:"podRef"`
		Starttime        time.Time          `json:"startTime"`
		Containers       []ContainerStruct  `json:"containers"`
		CPU              CPUStruct          `json:"cpu"`
		Memory           MemoryStruct       `json:"memory,omitempty"`
		Volume           []FsStruct         `json:"volume,omitempty"`
		EphemeralStorage FsStruct           `json:"ephemeral-storage"`
		ProcessStats     ProcessStatsStruct `json:"process_stats"`
	} `json:"pods"`
}

type MemoryStruct struct {
	Time            time.Time `json:"time"`
	Availablebytes  int64     `json:"availableBytes,omitempty"`
	Usagebytes      int64     `json:"usageBytes,omitempty"`
	Workingsetbytes int       `json:"workingSetBytes"`
	Rssbytes        int64     `json:"rssBytes,omitempty"`
	Pagefaults      int       `json:"pageFaults,omitempty"`
	Majorpagefaults int       `json:"majorPageFaults,omitempty"`
}

type CPUStruct struct {
	Time                 time.Time `json:"time"`
	Usagenanocores       int       `json:"usageNanoCores"`
	Usagecorenanoseconds int64     `json:"usageCoreNanoSeconds"`
}

type NetworkStruct struct {
	Time       time.Time `json:"time"`
	Name       string    `json:"name"`
	Rxbytes    int64     `json:"rxBytes"`
	Rxerrors   int       `json:"rxErrors"`
	Txbytes    int64     `json:"txBytes"`
	Txerrors   int       `json:"txErrors"`
	Interfaces []struct {
		Name     string `json:"name"`
		Rxbytes  int64  `json:"rxBytes"`
		Rxerrors int    `json:"rxErrors"`
		Txbytes  int64  `json:"txBytes"`
		Txerrors int    `json:"txErrors"`
	} `json:"interfaces"`
}

type FsStruct struct {
	Time           time.Time `json:"time"`
	Availablebytes int64     `json:"availableBytes,omitempty"`
	Capacitybytes  int64     `json:"capacityBytes"`
	Usedbytes      int64     `json:"usedBytes"`
	Inodesfree     int       `json:"inodesFree"`
	Inodes         int       `json:"inodes"`
	Inodesused     int       `json:"inodesUsed"`
	Name           string    `json:"name,omitempty"`
}

type RuntimeStruct struct {
	Imagefs struct {
		Time           time.Time `json:"time"`
		Availablebytes int64     `json:"availableBytes"`
		Capacitybytes  int64     `json:"capacityBytes"`
		Usedbytes      int64     `json:"usedBytes"`
		Inodesfree     int       `json:"inodesFree"`
		Inodes         int       `json:"inodes"`
		Inodesused     int       `json:"inodesUsed"`
	} `json:"imageFs"`
}

type RlimitStruct struct {
	Time    time.Time `json:"time"`
	Maxpid  int       `json:"maxpid"`
	Curproc int       `json:"curproc"`
}

type PodrefStruct struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
}

type ProcessStatsStruct struct {
	ProcessCount int `json:"process_count"`
}

type ContainerStruct struct {
	Name      string       `json:"name"`
	Starttime time.Time    `json:"startTime"`
	CPU       CPUStruct    `json:"cpu"`
	Memory    MemoryStruct `json:"memory,omitempty"`
	Rootfs    FsStruct     `json:"rootfs,omitempty"`
	Logs      FsStruct     `json:"logs,omitempty"`
}
