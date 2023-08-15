package snowflake

import (
	"errors"
	"sync"
	"time"
)

type Snowflake struct {
	mutex         sync.Mutex
	lastTimestamp int64
	workerID      int64
	sequence      int64
}

// NewSnowflake函数用于创建Snowflake实例

func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID >= 1024 {
		return nil, errors.New(" Worker ID must be between 0 and 1023")
	}

	return &Snowflake{
		lastTimestamp: -1,
		workerID:      workerID,
		sequence:      0,
	}, nil
}

// GenerateID生成唯一ID

func (sf *Snowflake) GenerateID() (int64, error) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	currentTimestamp := time.Now().UnixNano() / 1e6 // 毫秒级时间戳

	if currentTimestamp < sf.lastTimestamp {
		return -1, errors.New(" Clock moved backwards")
	}

	if currentTimestamp == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & 4095 // 序列号占12位，最大值为4095
		if sf.sequence == 0 {
			// 当前毫秒的序列号用完，等待下一毫秒
			for currentTimestamp <= sf.lastTimestamp {
				currentTimestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		sf.sequence = 0
	}

	sf.lastTimestamp = currentTimestamp

	// 生成ID：时间戳的42位 + Worker ID的10位 + 序列号的12位
	id := (currentTimestamp << 22) | (sf.workerID << 12) | sf.sequence

	return id, nil
}
