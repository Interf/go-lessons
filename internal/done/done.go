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

			<-ch

			once.Do(func() {
				close(orDone)
			})

		}(channel)
	}

	return orDone
}

func Or2(channels ...<-chan any) <-chan any {

	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan any)

	go func() {
		defer close(orDone)

		select {
		case <-channels[0]:
		case <-channels[1]:
		case <-Or2(channels[2:]...):
		}
	}()

	return orDone
}
