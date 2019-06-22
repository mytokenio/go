package metrics

import (
	"bytes"
	"encoding/json"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/mytokenio/go/log"
)

func cronMonitor() {
	t := time.NewTicker(cronInterval)

	for {
		select {
		case <-exitChan:
			return
		case <-t.C:
			reportStateFactory()
		}
	}
}

func reportStateFactory() {
	if globalKafka.producer == nil {
		return
	}

	value, err := getStateProducerValue()
	if err != nil {
		log.Errorf("getStateProducerValue err: %v", err)
		return
	}

	globalKafka.mutex.Lock()
	defer globalKafka.mutex.Unlock()

	if len(globalKafka.chanStateProducerValue) > default_producer_msg_caps-2 {
		log.Warnf("chanStateProducerValue is full, content: %s", value)
		return
	} else {
		globalKafka.chanStateProducerValue <- value
	}
}

func alarm(content string) {
	if content != "" && globalKafka.producer != nil {

		globalKafka.mutex.Lock()
		defer globalKafka.mutex.Unlock()

		if len(globalKafka.chanAlarmProducerValue) > default_producer_msg_caps-2 {
			log.Errorf("chanAlarmProducer is full, content: %s", content)
			return
		} else {
			ra := ReportAlarmPkg{
				JobID:       globalServiceInfo.jobID,
				ServiceName: globalServiceInfo.serviceName,
				Content:     content,
				HeartTime:   time.Now().Unix(),
			}
			pkg, _ := json.Marshal(ra)
			globalKafka.chanAlarmProducerValue <- string(pkg)
		}
	}
}

func reportMonitorCenter() {
	var value string
	var isNotClosed bool

	for {
		select {

		// receive value from chan_alarm
		case value, isNotClosed = <-globalKafka.chanAlarmProducerValue:
			if !isNotClosed {
				return
			}

			globalKafka.producer.Input() <- &sarama.ProducerMessage{
				Topic: globalKafka.reportAlarmTopic,
				Key:   sarama.StringEncoder(alarm_producter_msg_key),
				Value: sarama.ByteEncoder(value),
			}

		// receive value from chan_state
		case value, isNotClosed = <-globalKafka.chanStateProducerValue:
			if !isNotClosed {
				return
			}

			globalKafka.producer.Input() <- &sarama.ProducerMessage{
				Topic: globalKafka.reportStateTopic,
				Key:   sarama.StringEncoder(state_producter_msg_key),
				Value: sarama.ByteEncoder(value),
			}
		}
	}
}

func callback() error {
	var sucValue, failValue []byte
	var suc *sarama.ProducerMessage
	var fail *sarama.ProducerError

	for {
		select {
		case <-exitChan:
			return nil
		case suc = <-globalKafka.producer.Successes():
			sucValue, _ = suc.Value.Encode()
			if suc.Topic == globalKafka.reportAlarmTopic {
				log.Infof("send alarm msg success. [T:%s P:%d O:%d M:%s]",
					suc.Topic, suc.Partition, suc.Offset, string(sucValue))
			} else {
				log.Infof("send state msg success. [T:%s P:%d O:%d M:%s]",
					suc.Topic, suc.Partition, suc.Offset, string(sucValue))
			}

		case fail = <-globalKafka.producer.Errors():
			failValue, _ = fail.Msg.Value.Encode()
			if fail.Msg.Topic == globalKafka.reportAlarmTopic {
				log.Errorf("send alarm msg failed: [T:%s M:%s], err: %s",
					fail.Msg.Topic, string(failValue), fail.Err.Error())
			} else {
				log.Errorf("send state msg failed: [T:%s M:%s], err: %s",
					fail.Msg.Topic, string(failValue), fail.Err.Error())
			}

		}
	}
}

func getMemoryPercent() int {
	pid := strconv.Itoa(globalServiceInfo.processID)

	ps := exec.Command("ps", "-eo", "pid,pmem")
	grep := exec.Command("grep", pid)
	result, _, err := pipeline(ps, grep)
	if err != nil {
		log.Errorf("get process %mem err: %v", err)
		return 0
	}

	lines := strings.Split(string(result), "\n")
	for i := 0; i < len(lines); i++ {
		lineFields := strings.Fields(lines[i])
		if len(lineFields) == 2 && lineFields[0] == pid {
			men, _ := strconv.ParseFloat(lineFields[1], 10)
			memPer, _ := math.Modf(men)
			return int(memPer)
		}
	}

	return 0
}

func pipeline(cmds ...*exec.Cmd) ([]byte, []byte, error) {
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	var output bytes.Buffer
	var stderr bytes.Buffer
	var err error
	maxIndex := len(cmds) - 1
	cmds[maxIndex].Stdout = &output
	cmds[maxIndex].Stderr = &stderr

	for i, cmd := range cmds[:maxIndex] {
		if i == maxIndex {
			break
		}

		cmds[i+1].Stdin, err = cmd.StdoutPipe()
		if err != nil {
			return nil, nil, err
		}
	}

	// Start each command
	for _, cmd := range cmds {
		err := cmd.Start()
		if err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		err := cmd.Wait()
		if err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	return output.Bytes(), stderr.Bytes(), nil
}

func getStateProducerValue() (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	extend := make(map[string]interface{})
	now := time.Now().Unix()
	rs := ReportStatePkg{
		JobID:       globalServiceInfo.jobID,
		ServiceName: globalServiceInfo.serviceName,
		EnvType:     globalServiceInfo.envType,
		Host:        globalServiceInfo.host,
		ProcessID:   globalServiceInfo.processID,
		Memory:      getMemoryPercent(),
	}

	if v, ok := gaugeIntMap["status"]; !ok {
		rs.Status = STATUS_OK
	} else {
		rs.Status = int(v)
	}

	if v, ok := gaugeIntMap["start_time"]; ok {
		rs.StartTime = v
	}

	if v, ok := gaugeIntMap["heart_time"]; ok {
		rs.HeartTime = v
	} else {
		rs.HeartTime = now
	}

	if v, ok := gaugeIntMap["stop_time"]; ok {
		rs.StopTime = v
	}

	if v, ok := gaugeIntMap["exit_code"]; ok {
		rs.ExitCode = int(v)
		if v > 0 {
			rs.StopTime = now
		}
	}

	// set extend
	for key, value := range countMap {
		if !isDefaultKey[key] {
			extend[key] = value
		}
	}
	for key, value := range gaugeIntMap {
		if !isDefaultKey[key] {
			extend[key] = value
		}
	}
	for key, value := range gaugeStrMap {
		if !isDefaultKey[key] {
			extend[key] = value
		}
	}

	if len(extend) > 0 {
		if extendC, err := json.Marshal(extend); err == nil {
			rs.Extend = string(extendC)
		}
	}

	value, err := json.Marshal(rs)
	if err != nil {
		return "", err
	}

	return string(value), nil
}
