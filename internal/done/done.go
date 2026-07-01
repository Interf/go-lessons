package done

import (
	"sync"
)

func Or(channels ...<-chan any) <-chan any {
	if len(channels) == 0 {
		return nil
	}

	if len(channels) == 1 {
		return channels[0]
	}

	orDone := make(chan any)

	var once sync.Once

	for _, channel := range channels {
		go func(ch <-chan any) {

			_, ok := <-ch

			if !ok {
				once.Do(func() {
					close(orDone)
				})
			}

		}(channel)
	}

	return orDone
}
