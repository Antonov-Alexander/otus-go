package storages

import (
	"sync"
	"time"

	"github.com/Antonov-Alexander/otus-go/final_project/types"
)

// Используется алгоритм по мотивам Скользящего окна (Sliding Window)
// Делим интервал не более чем на 10 отрезков и не меньше 1-й секунды
// Высчитываем лимит запросов для одного отрезка
// кейс 1: 15 раз в 30 сек, расчёт: отрезок = 30 / 15 = 2 сек,   лимит = 1 запрос
// кейс 2: 30 раз в 15 сек, расчёт: отрезок = 15 / 10 = 1.5 сек, лимит = 30 / 10 = 3 запроса

const BucketsCount = 10
const Precision = 1000

type BucketInfo struct {
	Interval int
	Limit    int
	Count    int
}

type BucketValue struct {
	Timestamp int
	Value     int
}

type MemoryStorage struct {
	mux      *sync.Mutex
	info     map[types.Limit]BucketInfo
	counters map[types.Item][]BucketValue
	count    int
}

func (m *MemoryStorage) Init() {
	m.mux = &sync.Mutex{}
	m.info = map[types.Limit]BucketInfo{}
	m.counters = map[types.Item][]BucketValue{}
}

func (m *MemoryStorage) Inc(item types.Item, limit types.Limit) bool {
	m.mux.Lock()
	defer m.mux.Unlock()

	currentTimestamp := m.getCurrentTimestamp()
	bucketInfo := m.getBucketInfo(limit)

	// если нет каунтера, то создаём
	counter, ok := m.counters[item]
	if !ok {
		m.counters[item] = m.initCounter(bucketInfo, currentTimestamp)
		m.count++
		return true
	}

	// чистим каунтер и собираем инфу
	availableIdx := -1
	clearedIdx := -1
	var clearedCount int
	var maxActiveTimestamp int    // нужен чтобы создать следующий бакет после последнего созданного
	var minAvailableTimestamp int // нужен чтобы найти минимальный свободный бакет

	for idx, bucket := range counter {
		// удялаем протухший бакет
		if bucket.Timestamp < currentTimestamp-limit.Interval*Precision {
			clearedCount++
			counter[idx] = BucketValue{}
			clearedIdx = idx
		} else {
			// ищем максимальный timestamp среди бакетов
			if bucket.Timestamp > maxActiveTimestamp {
				maxActiveTimestamp = bucket.Timestamp
			}

			// ищем минимальный бакет, который можно заинкрементить
			if minAvailableTimestamp < bucket.Timestamp && bucket.Value < bucketInfo.Limit {
				minAvailableTimestamp = bucket.Timestamp
				availableIdx = idx
			}
		}
	}

	// если все бакеты удалили, то нарежем новый каунтер
	if clearedCount == len(counter) {
		m.counters[item] = m.initCounter(bucketInfo, currentTimestamp)
		return true
	}

	// если есть свободный бакет, то инкрементим его
	if availableIdx != -1 {
		m.counters[item][availableIdx].Value++
		return true
	}

	// если есть очищенный бакет, то инитим его
	if clearedIdx != -1 {
		m.counters[item][clearedIdx] = BucketValue{
			Timestamp: maxActiveTimestamp + bucketInfo.Interval,
			Value:     1,
		}

		return true
	}

	// не нашлось ненаполненных и очищенных бакетов
	return false
}

func (m *MemoryStorage) Reset(item types.Item) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if _, ok := m.counters[item]; ok {
		delete(m.counters, item)
	}
}

func (m *MemoryStorage) Cleanup() {
	m.mux.Lock()
	defer m.mux.Unlock()

	currentTimestamp := m.getCurrentTimestamp()
	for key, counter := range m.counters {
		var count int
		for _, bucket := range counter {
			if bucket.Timestamp < currentTimestamp {
				count++
			}
		}

		if count == len(counter) {
			delete(m.counters, key)
		}
	}
}

func (m *MemoryStorage) getCurrentTimestamp() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}

func (m *MemoryStorage) getBucketInfo(limit types.Limit) BucketInfo {
	res, ok := m.info[limit]
	if !ok {
		res = m.calcBucket(limit)
		m.info[limit] = res
	}

	return res
}

func (*MemoryStorage) calcBucket(limit types.Limit) (bucket BucketInfo) {
	if limit.Limit <= limit.Interval {
		// Кейс 1: лимит < интервала
		bucket.Interval = Precision * limit.Interval / limit.Limit
		bucket.Limit = 1
	} else {
		bucketsCount := BucketsCount
		if bucketsCount > limit.Limit {
			bucketsCount = limit.Limit
		}

		// Кейс 2: лимит > интервала
		bucket.Interval = Precision * limit.Interval / bucketsCount
		bucket.Limit = limit.Limit / bucketsCount
	}

	bucket.Count = Precision * limit.Interval / bucket.Interval
	return
}

func (m *MemoryStorage) initCounter(bucketInfo BucketInfo, startTimestamp int) []BucketValue {
	result := make([]BucketValue, bucketInfo.Count)
	for idx := range result {
		result[idx].Timestamp = startTimestamp + bucketInfo.Interval*idx
	}

	result[0].Value = 1
	return result
}
