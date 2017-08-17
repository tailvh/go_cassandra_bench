/*
 */
package logger

import (
  "fmt"
  "os"
  "os/user"
  "path"
  "runtime"
  "sort"
  "strconv"
  "strings"
  "sync"
  "time"
  "zmqlog/rxconfig"
)

// consts
const (
  kMaxInt64          = int64(^uint64(0) >> 1)
  kLogCreatedTimeLen = 24
  kLogFilenameMinLen = 29
)

// log level
const (
  kLogLevelTrace = iota
  kLogLevelInfo
  kLogLevelWarn
  kLogLevelError
  kLogLevelPanic
  kLogLevelAbort
  kLogLevelMax
)

// log flags
const (
  kFlagLogTrace = 1 << iota
  kFlagLogThrough
  kFlagLogFuncName
  kFlagLogFilenameLineNum
  kFlagLogToConsole
)

// const strings
const (
  // Default filename prefix for logfiles
  DefFilenamePrefix = "LOG_CHARGE_CARD"
  // Default filename prefix for symlinks to logfiles
  DefSymlinkPrefix = "LOG_CHARGE_CARD_CURRENT"
  kLogLevelChar    = "TIWEPA"
)

// Init must be called first, otherwise this logger will not function properly!
// It returns nil if all goes well, otherwise it returns the corresponding error.
//   maxfiles: Must be greater than 0 and less than or equal to 100000.
//   nfilesToDel: Number of files deleted when number of log files reaches `maxfiles`.
//                Must be greater than 0 and less than or equal to `maxfiles`.
//   maxsize: Maximum size of a log file in MB, 0 means unlimited.
//   logTrace: If set to false, `logger.Trace("xxxx")` will be mute.
func Init(logpath string, maxfiles, nfilesToDel int, maxsize uint32, logTrace bool) error {
  err := os.MkdirAll(logpath, 0755)
  if err != nil {
    return err
  }
  if maxfiles <= 0 || maxfiles > 100000 {
    return fmt.Errorf("maxfiles must be greater than 0 and less than or equal to 100000: %d", maxfiles)
  }
  if nfilesToDel <= 0 || nfilesToDel > maxfiles {
    return fmt.Errorf("nfilesToDel must be greater than 0 and less than or equal to maxfiles! toDel=%d maxfiles=%d",
      nfilesToDel, maxfiles)
  }
  gConf.logPath = logpath + "/"
  gConf.setFlags(kFlagLogTrace, logTrace)
  gConf.maxfiles = maxfiles
  gConf.nfilesToDel = nfilesToDel
  gConf.setMaxSize(maxsize)
  return SetFilenamePrefix(DefFilenamePrefix, DefSymlinkPrefix)
}

// SetLogThrough sets whether to write log to all the logfiles with less severe log level.
// By default, logthrough is turn on. You can turn it off for better performance.
func SetLogThrough(on bool) {
  gConf.setFlags(kFlagLogThrough, on)
}

// SetLogFunctionName sets whether to log down the function name where the log takes place.
// By default, function name is not logged down for better performance.
func SetLogFunctionName(on bool) {
  gConf.setFlags(kFlagLogFuncName, on)
}

// SetLogFilenameLineNum sets whether to log down the filename and line number where the log takes place.
// By default, filename and line number are logged down. You can turn it off for better performance.
func SetLogFilenameLineNum(on bool) {
  gConf.setFlags(kFlagLogFilenameLineNum, on)
}

// SetLogToConsole sets whether to output logs to the console.
// By default, logs are not output to the console.
func SetLogToConsole(on bool) {
  gConf.setFlags(kFlagLogToConsole, on)
}

// SetFilenamePrefix sets filename prefix for the logfiles and symlinks of the logfiles.
//
// Filename format for logfiles is `PREFIX`.`SEVERITY_LEVEL`.`DATE_TIME`.log
//
// Filename format for symlinks is `PREFIX`.`SEVERITY_LEVEL`
//
// 3 kinds of placeholders can be used in the prefix: %P, %H and %U.
//
// %P means program name, %H means hostname, %U means username.
//
// The default prefix for a log filename is logger.DefFilenamePrefix ("%P.%H.%U").
// The default prefix for a symlink is logger.DefSymlinkPrefix ("%P.%U").
func SetFilenamePrefix(logfilenamePrefix, symlinkPrefix string) error {
  gConf.setFilenamePrefix(logfilenamePrefix, symlinkPrefix)
  files, err := getLogfilenames(gConf.logPath)
  if err == nil {
    gConf.curfiles = len(files)
  }
  return err
}

// Trace logs down a log with trace level.
// If parameter logTrace of logger.Init() is set to be false, no trace logs will be logged down.
func Trace(format string, args ...interface{}) {
  if gConf.logTrace() {
    log(kLogLevelTrace, format, args)
  }
}

// Info logs down a log with info level.
func Info(format string, args ...interface{}) {
  log(kLogLevelInfo, format, args)
}

// Warn logs down a log with warning level.
func Warn(format string, args ...interface{}) {
  log(kLogLevelWarn, format, args)
}

// Error logs down a log with error level.
func Error(format string, args ...interface{}) {
  log(kLogLevelError, format, args)
}

// Panic logs down a log with panic level and then panic("panic log") is called.
func Panic(format string, args ...interface{}) {
  log(kLogLevelPanic, format, args)
  panic("panic log")
}

// Abort logs down a log with abort level and then os.Exit(-1) is called.
func Abort(format string, args ...interface{}) {
  log(kLogLevelAbort, format, args)
  os.Exit(-1)
}

// logger configuration
type config struct {
  logPath     string
  pathPrefix  string
  logflags    uint32
  maxfiles    int   // limit the number of log files under `logPath`
  curfiles    int   // number of files under `logPath` currently
  nfilesToDel int   // number of files deleted when reaching the limit of the number of log files
  maxsize     int64 // limit size of a log file
  purgeLock   sync.Mutex
}

func (conf *config) setFlags(flag uint32, on bool) {
  if on {
    conf.logflags = conf.logflags | flag
  } else {
    conf.logflags = conf.logflags & ^flag
  }
}
func (conf *config) logTrace() bool {
  return (conf.logflags & kFlagLogTrace) != 0
}
func (conf *config) logThrough() bool {
  return (conf.logflags & kFlagLogThrough) != 0
}
func (conf *config) logFuncName() bool {
  return (conf.logflags & kFlagLogFuncName) != 0
}
func (conf *config) logFilenameLineNum() bool {
  return (conf.logflags & kFlagLogFilenameLineNum) != 0
}
func (conf *config) logToConsole() bool {
  return (conf.logflags & kFlagLogToConsole) != 0
}
func (conf *config) setMaxSize(maxsize uint32) {
  if maxsize > 0 {
    conf.maxsize = int64(maxsize) * 1024 * 1024
  } else {
    conf.maxsize = kMaxInt64 - (1024 * 1024 * 1024 * 1024 * 1024)
  }
}
func (conf *config) setFilenamePrefix(filenamePrefix, symlinkPrefix string) {
  host, err := os.Hostname()
  if err != nil {
    host = "Unknown"
  }
  username := "Unknown"
  curUser, err := user.Current()
  if err == nil {
    tmpUsername := strings.Split(curUser.Username, "\\") // for compatible with Windows
    username = tmpUsername[len(tmpUsername)-1]
  }
  conf.pathPrefix = conf.logPath
  if len(filenamePrefix) > 0 {
    filenamePrefix = strings.Replace(filenamePrefix, "%P", gProgname, -1)
    filenamePrefix = strings.Replace(filenamePrefix, "%H", host, -1)
    filenamePrefix = strings.Replace(filenamePrefix, "%U", username, -1)
    //filenamePrefix = strings.Replace(filenamePrefix, "%D", , -1)
    conf.pathPrefix = conf.pathPrefix + filenamePrefix // + "."
  }
  if len(symlinkPrefix) > 0 {
    symlinkPrefix = strings.Replace(symlinkPrefix, "%P", gProgname, -1)
    symlinkPrefix = strings.Replace(symlinkPrefix, "%H", host, -1)
    symlinkPrefix = strings.Replace(symlinkPrefix, "%U", username, -1)
    //symlinkPrefix += "."
  }
  isSymlink = map[string]bool{}
  for i := 0; i != kLogLevelMax; i++ {
    gLoggers[i].level = i
    gSymlinks[i] = symlinkPrefix //+ gLogLevelNames[i]
    isSymlink[gSymlinks[i]] = true
    gFullSymlinks[i] = conf.logPath + gSymlinks[i]
  }
}

// logger
type logger struct {
  file  *os.File
  level int
  day   int
  size  int64
  lock  sync.Mutex
}

func (l *logger) log(t time.Time, data []byte) {
  y, m, d := t.Date()
  l.lock.Lock()
  defer l.lock.Unlock()
  filestat, _ := l.file.Stat()
  if filestat != nil {
    l.size = filestat.Size()
  } else {
    // Just restart log server
    files, _ := getLogfilenames(gConf.logPath)
    sort.Sort(byCreatedTime(files))
    // This when server log is turn off and need tracking what file to process
    curfilename := "rxnofile"
    for _, element := range files {
      tempdatestr := fmt.Sprintf("%d-%02d-%02d-", y, m, d)
      if strings.Contains(element, tempdatestr) {
        curfilename = element
      }
    }
    if len(files) > 0 && curfilename != "rxnofile" {
      tempfile, _ := os.OpenFile(gConf.logPath+curfilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
      defer tempfile.Close()
      filestat, _ = tempfile.Stat()
      l.size = filestat.Size()
      l.file = tempfile
      l.day = d
    }
  }
  if l.size >= gConf.maxsize || l.day != d || l.file == nil {
    // When new file, difference day, oversize
    // if l.size >= 5000 || l.day != d || l.file == nil {
    //hour, min, sec := t.Clock()
    filenumstr := "00000"
    // If full file need rotate +1 number file
    if l.day == d {
      tempFilename := filestat.Name()
      last5 := tempFilename[len(tempFilename)-5:]
      tempnum, _ := strconv.Atoi(last5)
      filenumstr = fmt.Sprintf("%05d", tempnum+1)
    }
    gConf.purgeLock.Lock()
    hasLocked := true
    defer func() {
      if hasLocked {
        gConf.purgeLock.Unlock()
      }
    }()
    filename := fmt.Sprintf("%s-%d-%02d-%02d-"+filenumstr, gConf.pathPrefix, y, m, d)
    newfile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
      l.errlog(t, data, err)
      return
    }
    gConf.curfiles++
    gConf.purgeLock.Unlock()
    hasLocked = false
    l.file.Close()
    l.file = newfile
    l.day = d
    l.size = 0
    err = os.RemoveAll(gFullSymlinks[l.level])
    if err != nil {
      l.errlog(t, nil, err)
    }
    err = os.Symlink(path.Base(filename), gFullSymlinks[l.level])
    if err != nil {
      l.errlog(t, nil, err)
    }
  }
  //n, err := l.file.Write(data)
  _, _ = l.file.Write(data)
  if filestat != nil {
    l.size = filestat.Size()
  } else {
    l.size = 0
  }
}

// (l *logger).errlog() should only be used within (l *logger).log()
func (l *logger) errlog(t time.Time, originLog []byte, err error) {
  buf := gBufPool.getBuffer()
  genLogPrefix(buf, l.level, 2, t)
  buf.WriteString(err.Error())
  buf.WriteByte('\n')
  if l.file != nil {
    l.file.Write(buf.Bytes())
    if len(originLog) > 0 {
      l.file.Write(originLog)
    }
  } else {
    fmt.Fprint(os.Stderr, buf.String())
    if len(originLog) > 0 {
      fmt.Fprint(os.Stderr, string(originLog))
    }
  }
  gBufPool.putBuffer(buf)
}

// sort files by created time embedded in the filename
type byCreatedTime []string

func (a byCreatedTime) Len() int {
  return len(a)
}
func (a byCreatedTime) Less(i, j int) bool {
  s1, s2 := a[i], a[j]
  if len(s1) < kLogFilenameMinLen {
    return true
  } else if len(s2) < kLogFilenameMinLen {
    return false
  } else {
    return s1[len(s1)-kLogCreatedTimeLen:] < s2[len(s2)-kLogCreatedTimeLen:]
  }
}
func (a byCreatedTime) Swap(i, j int) {
  a[i], a[j] = a[j], a[i]
}

// init is called after all the variable declarations in the package have evaluated their initializers,
// and those are evaluated only after all the imported packages have been initialized.
// Besides initializations that cannot be expressed as declarations, a common use of init functions is to verify
// or repair correctness of the program state before real execution begins.
func init() {
  tmpProgname := strings.Split(gProgname, "\\") // for compatible with `go run` under Windows
  gProgname = tmpProgname[len(tmpProgname)-1]
  gConf.setFilenamePrefix(DefFilenamePrefix, DefSymlinkPrefix)
}

// helpers
func getLogfilenames(dir string) ([]string, error) {
  var filenames []string
  f, err := os.Open(dir)
  if err == nil {
    filenames, err = f.Readdirnames(0)
    f.Close()
    if err == nil {
      nfiles := len(filenames)
      for i := 0; i < nfiles; {
        if isSymlink[filenames[i]] == false {
          i++
        } else {
          nfiles--
          filenames[i] = filenames[nfiles]
          filenames = filenames[:nfiles]
        }
      }
    }
  }
  return filenames, err
}
func genLogPrefix(buf *buffer, logLevel, skip int, t time.Time) {
  h, m, s := t.Clock()
  // time
  buf.tmp[0] = kLogLevelChar[logLevel]
  buf.twoDigits(1, h)
  buf.tmp[3] = ':'
  buf.twoDigits(4, m)
  buf.tmp[6] = ':'
  buf.twoDigits(7, s)
  buf.Write(buf.tmp[:9])
  var pc uintptr
  var ok bool
  if gConf.logFilenameLineNum() {
    var file string
    var line int
    pc, file, line, ok = runtime.Caller(skip)
    if ok {
      buf.WriteByte(' ')
      buf.WriteString(path.Base(file))
      buf.tmp[0] = ':'
      n := buf.someDigits(1, line)
      buf.Write(buf.tmp[:n+1])
    }
  }
  if gConf.logFuncName() {
    if !ok {
      pc, _, _, ok = runtime.Caller(skip)
    }
    if ok {
      buf.WriteByte(' ')
      buf.WriteString(runtime.FuncForPC(pc).Name())
    }
  }
  buf.WriteString("] ")
}
func log(logLevel int, format string, args []interface{}) {
  buf := gBufPool.getBuffer()
  t := time.Now()
  // genLogPrefix(buf, logLevel, 3, t)
  fmt.Fprintf(buf, format, args...)
  buf.WriteByte('\n')
  output := buf.Bytes()
  if gConf.logThrough() {
    for i := logLevel; i != kLogLevelTrace; i-- {
      gLoggers[i].log(t, output)
    }
    if gConf.logTrace() {
      gLoggers[kLogLevelTrace].log(t, output)
    }
  } else {
    gLoggers[logLevel].log(t, output)
  }
  if gConf.logToConsole() {
    fmt.Print(string(output))
  }
  gBufPool.putBuffer(buf)
}

var gProgname = path.Base(os.Args[0])
var gLogLevelNames = [kLogLevelMax]string{
  "TRACE", "INFO", "WARN", "ERROR", "PANIC", "ABORT",
}
var gConf = config{
  logPath:     rxconfig.ScribeDirectory,
  logflags:    kFlagLogFilenameLineNum | kFlagLogThrough,
  maxfiles:    1000000000,
  nfilesToDel: 1,
  maxsize:     10 * 1024 * 1024,
}
var gSymlinks [kLogLevelMax]string
var isSymlink map[string]bool
var gFullSymlinks [kLogLevelMax]string
var gBufPool bufferPool
var gLoggers [kLogLevelMax]logger
