package wrapper

import (
	// #cgo pkg-config: swisseph
	// #include <stdio.h>
	// #include <errno.h>
	// #include "swephexp.h"
	"C"
)
import (
	"math"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/afjoseph/sacredstar/projectpath"
)

type SwissEph struct{}

func NewWithBuiltinPath() *SwissEph {
	return NewWithPath(
		filepath.Join(projectpath.Root, "ext", "swisseph"),
	)
}

func NewWithPath(path string) *SwissEph {
	path = filepath.Join(path, "ephe")

	C.swe_set_ephe_path(C.CString(path))
	C.swe_set_sid_mode(C.SE_SIDM_LAHIRI, 0, 0)
	return &SwissEph{}
}

func (s *SwissEph) Close() {
	C.swe_close()
}

func (s *SwissEph) Version() string {
	b := make([]byte, 256)
	a := (*C.char)(unsafe.Pointer(C.CBytes(b)))
	defer C.free(unsafe.Pointer(a))
	C.swe_version(a)
	return C.GoString(a)
}

func (s *SwissEph) GoTimeToJulianDay(date time.Time) float64 {
	// XXX <02-02-2024, afjoseph> swisseph doesn't work with anything lower
	// than 'hours', but it takes hours as a double, so you can add minutes
	// to it as fractions. We don't work with seconds here
	hoursAsFraction := float64(date.Hour()) +
		float64(date.Minute())/60
	julDay := C.swe_julday(
		C.int(date.Year()),
		C.int(date.Month()),
		C.int(date.Day()),
		C.double(float64(hoursAsFraction)),
		C.int(1), // Gregorian calendar
	)
	return float64(julDay)
}

func (s *SwissEph) JulianDayToGoTime(julDay float64) time.Time {
	year := make([]C.int, 1)
	yearPtr := (*C.int)(unsafe.Pointer(&year[0]))
	month := make([]C.int, 1)
	monthPtr := (*C.int)(unsafe.Pointer(&month[0]))
	day := make([]C.int, 1)
	dayPtr := (*C.int)(unsafe.Pointer(&day[0]))
	hour := make([]C.double, 1)
	hourPtr := (*C.double)(unsafe.Pointer(&hour[0]))

	C.swe_revjul(
		C.double(julDay),
		C.int(1), // Gregorian calendar
		yearPtr,
		monthPtr,
		dayPtr,
		hourPtr,
	)

	// Calculate minutes from the fraction of the hour
	hoursAsInt := int(hour[0])
	hoursAsFraction := float64(hour[0]) - float64(hoursAsInt)
	minutes := math.Round(hoursAsFraction * 60)
	return time.Date(
		int(year[0]),
		time.Month(month[0]),
		int(day[0]),
		hoursAsInt, int(minutes), 0,
		0,
		time.UTC,
	)
}
