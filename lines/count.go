package lines

import "io"

// CountFromReader counts lines from a reader
func CountFromReader(r io.Reader) (count int64, err error) {
	lines, errors := LiveCountFromReader(r)

	for partialCount := range lines {
		count = partialCount
	}

	err = <-errors
	return
}

// LiveCountFromReader counts lines from a reader with "live" updates
func LiveCountFromReader(r io.Reader) (<-chan int64, <-chan error) {
	lines := make(chan int64, 16)
	errors := make(chan error, 1)

	go func() {
		defer func() {
			close(lines)
			close(errors)
		}()

		buf := make([]byte, 4096)

		var count int64

		for {
			n, err := r.Read(buf)
			for i := 0; i < n; i++ {
				if buf[i] == '\n' {
					count++
				}
			}

			if err != nil {
				if err == io.EOF {
					// FIXME Should we count a newline here?
					err = nil
				}
				errors <- err
				return
			}
		}
	}()

	return lines, errors
}
