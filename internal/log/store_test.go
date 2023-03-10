package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	write          = []byte("hello world")
	width          = uint64(len(write)) + lenWidth
	numberOfCycles = uint64(4)
)

func TestStoreAppendRead(t *testing.T) {
	f, err := ioutil.TempFile("", "store_append_read_test")
	require.NoError(t, err)
	defer func(name string) {
		_ = os.Remove(name)
	}(f.Name())

	s, err := newStore(f) // creating new store
	require.NoError(t, err)

	testAppend(t, s) // checking append
	testRead(t, s)   // checking read
	testReadAt(t, s) // checking read by offset

	s, err = newStore(f) // create new store with the same temp file
	require.NoError(t, err)
	testRead(t, s) // is service recovered after restart? Checking here.
}

func testAppend(t *testing.T, s *store) {
	t.Helper()
	for i := uint64(1); i < numberOfCycles; i++ {
		n, pos, err := s.Append(write)
		require.NoError(t, err)
		require.Equal(t, pos+n, width*i)
	}
}

func testRead(t *testing.T, s *store) {
	t.Helper()
	var pos uint64
	for i := uint64(1); i < numberOfCycles; i++ {
		read, err := s.Read(pos)
		require.NoError(t, err)
		require.Equal(t, write, read)
		pos += width
	}
}

func testReadAt(t *testing.T, s *store) {
	t.Helper()
	for i, off := uint64(1), int64(0); i < numberOfCycles; i++ {
		b := make([]byte, lenWidth)
		n, err := s.ReadAt(b, off)
		require.NoError(t, err)
		require.Equal(t, lenWidth, n)
		off += int64(n)

		size := enc.Uint64(b)
		b = make([]byte, size)
		n, err = s.ReadAt(b, off)
		require.NoError(t, err)
		require.Equal(t, write, b)
		require.Equal(t, int(size), n)
		off += int64(n)
	}
}

func TestStoreClose(t *testing.T) {
	f, err := ioutil.TempFile("", "store_close_test_")
	require.NoError(t, err)
	defer func(name string) {
		_ = os.Remove(name)
	}(f.Name())

	s, err := newStore(f)
	require.NoError(t, err)
	_, _, err = s.Append(write)
	require.NoError(t, err)

	f, beforeSize, err := openFile(f.Name())
	require.NoError(t, err)

	err = s.Close()
	require.NoError(t, err)

	_, afterSize, err := openFile(f.Name())
	require.NoError(t, err)
	require.True(t, afterSize > beforeSize)
}

func openFile(name string) (file *os.File, size int64, err error) {
	f, err := os.OpenFile(
		name,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, 0, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}
	return f, fi.Size(), nil
}
