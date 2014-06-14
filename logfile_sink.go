package health

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"time"
)

// This sink writes bytes in a format that a human might like to read in a logfile
type LogfileWriterSink struct {
	Writer io.Writer
}

func (s *LogfileWriterSink) EmitEvent(job string, event string, kvs map[string]string) error {
	var b bytes.Buffer
	b.WriteRune('[')
	b.WriteString(timestamp())
	b.WriteString("]: job:")
	b.WriteString(job)
	b.WriteString(" event:")
	b.WriteString(event)
	writeMapConsistently(&b, kvs)
	b.WriteRune('\n')
	_, err := s.Writer.Write(b.Bytes())
	return err
}

func (s *LogfileWriterSink) EmitEventErr(job string, event string, inputErr error, kvs map[string]string) error {
	var b bytes.Buffer
	b.WriteRune('[')
	b.WriteString(timestamp())
	b.WriteString("]: job:")
	b.WriteString(job)
	b.WriteString(" event:")
	b.WriteString(event)
	b.WriteString(" err:")
	b.WriteString(inputErr.Error())
	writeMapConsistently(&b, kvs)
	b.WriteRune('\n')
	_, err := s.Writer.Write(b.Bytes())
	return err
}

func (s *LogfileWriterSink) EmitTiming(job string, event string, nanos int64, kvs map[string]string) error {
	var b bytes.Buffer
	b.WriteRune('[')
	b.WriteString(timestamp())
	b.WriteString("]: job:")
	b.WriteString(job)
	b.WriteString(" event:")
	b.WriteString(event)
	b.WriteString(" time:")
	writeNanoseconds(&b, nanos)
	writeMapConsistently(&b, kvs)
	b.WriteRune('\n')
	_, err := s.Writer.Write(b.Bytes())
	return err
}

func (s *LogfileWriterSink) EmitJobCompletion(job string, kind CompletionType, nanos int64, kvs map[string]string) error {
	var b bytes.Buffer
	b.WriteRune('[')
	b.WriteString(timestamp())
	b.WriteString("]: job:")
	b.WriteString(job)
	b.WriteString(" status:")
	b.WriteString(completionTypeToString[kind])
	b.WriteString(" time:")
	writeNanoseconds(&b, nanos)
	writeMapConsistently(&b, kvs)
	b.WriteRune('\n')
	_, err := s.Writer.Write(b.Bytes())
	return err
}

func timestamp() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func writeMapConsistently(b *bytes.Buffer, kvs map[string]string) {
	if kvs == nil {
		return
	}
	keys := make([]string, 0, len(kvs))
	for k := range kvs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	keysLenMinusOne := len(keys) - 1

	b.WriteString(" kvs:[")
	for i, k := range keys {
		b.WriteString(k)
		b.WriteRune(':')
		b.WriteString(kvs[k])

		if i != keysLenMinusOne {
			b.WriteRune(' ')
		}
	}
	b.WriteRune(']')
}

func writeNanoseconds(b *bytes.Buffer, nanos int64) {
	switch {
	case nanos > 2000000:
		fmt.Fprintf(b, "%d ms", nanos/1000000)
	case nanos > 2000:
		fmt.Fprintf(b, "%d μs", nanos/1000)
	default:
		fmt.Fprintf(b, "%d ns", nanos)
	}
}