package utils

import (
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	s := `
// ToMapStr returns a new MapStr containing the data from this Record.
func (e Record) ToEvent() beat.Event {
	m := common.MapStr{
		"type":          e.API,
		"log_name":      e.Channel,
		"source_name":   e.Provider.Name,
		"computer_name": e.Computer,
		"record_number": strconv.FormatUint(e.RecordID, 10),
		"event_id":      e.EventIdentifier.ID,
	}

	addOptional(m, "xml", e.XML)
	addOptional(m, "provider_guid", e.Provider.GUID)
	addOptional(m, "version", e.Version)
	addOptional(m, "level", e.Level)
	addOptional(m, "task", e.Task)
	addOptional(m, "opcode", e.Opcode)
	addOptional(m, "keywords", e.Keywords)
	addOptional(m, "message", sys.RemoveWindowsLineEndings(e.Message))
	addOptional(m, "message_error", e.RenderErr)

	// Correlation
	addOptional(m, "activity_id", e.Correlation.ActivityID)
	addOptional(m, "related_activity_id", e.Correlation.RelatedActivityID)

	// Execution
	addOptional(m, "process_id", e.Execution.ProcessID)
	addOptional(m, "thread_id", e.Execution.ThreadID)
	addOptional(m, "processor_id", e.Execution.ProcessorID)
	addOptional(m, "session_id", e.Execution.SessionID)
	addOptional(m, "kernel_time", e.Execution.KernelTime)
	addOptional(m, "user_time", e.Execution.UserTime)
	addOptional(m, "processor_time", e.Execution.ProcessorTime)

	if e.User.Identifier != "" {
		user := common.MapStr{
			"identifier": e.User.Identifier,
		}
		m["user"] = user

		addOptional(user, "name", e.User.Name)
		addOptional(user, "domain", e.User.Domain)
		addOptional(user, "type", e.User.Type.String())
	}

	addPairs(m, "event_data", e.EventData.Pairs)
	userData := addPairs(m, "user_data", e.UserData.Pairs)
	addOptional(userData, "xml_name", e.UserData.Name.Local)

	return beat.Event{
		Timestamp: e.TimeCreated.SystemTime,
		Fields:    m
		Private:   e.Offset,
	}
`
	t.Log(MD5(s))
}

func TestExcludeFile(t *testing.T) {
	p := "/sdf.swpb"
	t.Log(ExcludeFile(p))

	p = "/sdf.swp"
	t.Log(ExcludeFile(p))
}

func TestRemoveFileRotationSuffix(t *testing.T) {
	p := "/Users/liuh/development/temp/logs/another.log.2018-04-1"
	t.Log(RemoveFileRotationSuffix(p))
	p = "/sdf.swpb"
	t.Log(RemoveFileRotationSuffix(p))
	p = "swpb"
	t.Log(RemoveFileRotationSuffix(p))

}

func TestExtractTaskKeyFromPath(t *testing.T) {
	root := "/Users/liuh/development/temp/logs"
	tp := "/Users/liuh/development/temp/logs/app/sd.f.s.df"

	k, e := ExtractTaskKeyFromPath(tp, root, "app", "", "")
	if e != nil {
		t.Error(e)
	} else {
		t.Log("extract task key: ", k)
	}
}

func TestFormatTaskKey(t *testing.T) {
	ks := []string{
		"sdf-sdf-sdf",
		"sdf.sdf.sdf",
		"sdf3-SDF-sdf",
		"sdf-sdf!-sdf",
		"sdf-sdf-|@sdf",
	}

	for _, k := range ks {
		fk, err := FormatTaskKey(k)
		if err != nil {
			t.Log(k, " format error!", err)
		} else {
			t.Log(k, " format success! ", fk)
		}
	}
}

func TestCurOffset(t *testing.T) {
	p := "/Users/liuh/development/temp/tcpdump/dns.84.dump"
	f, err := ReadOpen(p)
	if err != nil {
		t.Error(err)
	}

	s, err := CurOffset(f)
	if err != nil {
		t.Error(err)
	}

	t.Log("seek end offset: ", s)

	st, _ := f.Stat()
	t.Log("file size: ", st.Size())
}

func TestCheckFileLogType(t *testing.T) {
	file,_,err1:=GetFileStat("/home/tqyin")
	fmt.Println(file.IsDir(),err1)
	p1 := "/Users/liuh/development/temp/go.log"
	logType, err := CheckFileLogType(p1)
	if err != nil {
		t.Error(err)
	}
	t.Log(p1, "log type =", logType)

	p2 := "/Users/liuh/development/temp/go.log.20180425"
	logType, err = CheckFileLogType(p2)
	if err != nil {
		t.Error(err)
	}
	t.Log(p1, "log type =", logType)

	p3 := "/Users/liuh/development/temp/go.json"
	logType, err = CheckFileLogType(p3)
	if err != nil {
		t.Error(err)
	}
	t.Log(p1, "log type =", logType)

	p4 := "/Users/liuh/development/temp/go.json.2018-04-25"
	logType, err = CheckFileLogType(p4)
	if err != nil {
		t.Error(err)
	}
	t.Log(p1, "log type =", logType)
}

func TestIsLogFile(t *testing.T) {
	f := "/Users/liuh/development/logs/platform.log"
	if !IsLogFile(f) {
		t.Error("failed")
	}

	f = "/Users/liuh/development/logs/platform.1log"
	if IsLogFile(f) {
		t.Error("failed")
	}

	f = "/Users/liuh/development/logs/platform.json"
	if !IsLogFile(f) {
		t.Error("failed")
	}

	f = "/Users/liuh/development/logs/platform.js1on"
	if IsLogFile(f) {
		t.Error("failed")
	}
}
