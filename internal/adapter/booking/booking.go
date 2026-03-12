package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/cybernetlab/swimming-search/internal/domain"
)

type Config struct {
	BaseURL string `envconfig:"BOOKING_BASE_URL" required:"true"`
}

type Booking struct {
	config Config
}

type CentreResult struct {
	CourseGroupCategoryIds []int
	Centres                []domain.Centre `json:"value"`
	UserData               any             `json:"-"`
}

type SearchResult struct {
	Error          string
	Message        string
	FiltersApplied domain.Filter
	ResultSet      struct {
		GroupedBy string
		Results   []domain.Course
	}
}

func New(config Config) *Booking {
	return &Booking{config: config}
}

func (b *Booking) GetCentres(ctx context.Context) ([]domain.Centre, error) {
	url := fmt.Sprintf("%s/findLocation", b.config.BaseURL)

	ctx, cancel := context.WithCancelCause(ctx)
	timer := time.AfterFunc(5*time.Second, func() { cancel(context.DeadlineExceeded) })
	defer timer.Stop()
	out := make(chan []domain.Centre)

	go func() {
		result, err := fetch[CentreResult](url)
		if err != nil {
			cancel(fmt.Errorf("fetch[CentreResult](%s): %w", url, err))
		}
		out <- result.Centres
	}()

	select {
	case centres := <-out:
		return centres, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (b *Booking) StartSearchCourses(ctx context.Context, search domain.Search, out chan<- domain.Course) {
	go func() {
		ctx, cancel := context.WithCancel(ctx)

		var wg sync.WaitGroup

		for _, centre := range search.CentreIDs {
			wg.Go(func() { searchCourses(ctx, b.config.BaseURL, centre, search, out, cancel) })
		}

		<-ctx.Done()

		cancel()
		wg.Wait()
		close(out)
	}()
}

func searchCourses(
	ctx context.Context,
	baseURL string,
	centre uint,
	search domain.Search,
	out chan<- domain.Course,
	cancel context.CancelFunc,
) {
	for {
		courses := []domain.Course{}
		filter := domain.Filter{
			Region:              1,
			ShowFullCourses:     false,
			Centre:              centre,
			CourseGroupCategory: []uint{1},
			Offset:              0,
			Limit:               20,
		}

		err := fetchCourses(baseURL, filter, &courses)
		if err == nil {
			found := 0
			for _, course := range courses {
				if course.IsMatch(search) {
					out <- course
					found++
				}
			}
			if found > 0 {
				log.Info().Msgf("%d courses found for centre %d. Search stopped", found, centre)
				cancel()
				return
			} else {
				log.Info().Msgf("No courses found for centre %d. Next search in 1 hour", centre)
			}
			select {
			case <-ctx.Done():
				log.Debug().Msgf("Search for centre %d is closed", centre)
				return
			case <-time.After(time.Hour):
				continue
			}
		} else {
			log.Error().Err(err).Msgf("Error fetching courses for centre %d", centre)
			select {
			case <-ctx.Done():
				log.Debug().Msgf("Search for centre %d is closed", centre)
				return
			case <-time.After(5 * time.Minute):
				continue
			}
		}
	}
}

func fetch[T SearchResult | CentreResult](url string) (T, error) {
	var result T

	log.Debug().Msgf("Making HTTP GET request to: %s", url)
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, fmt.Errorf("http.NewRequest: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("io.ReadAll: %w", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return result, nil
}

func fetchCourses(baseURL string, filter domain.Filter, courses *[]domain.Course) error {
	q, _ := json.Marshal(filter)
	url := fmt.Sprintf("%s/search?filter=%s", baseURL, url.QueryEscape(string(q)))

	result, err := fetch[SearchResult](url)
	if err != nil {
		return fmt.Errorf("fetch[SearchResult](%s): %w", url, err)
	}
	switch {
	case result.Error != "":
		return fmt.Errorf("%w: %s", domain.ErrAPI, result.Error)
	case result.Message != "":
		return fmt.Errorf("%w: %s", domain.ErrAPI, result.Message)
	}

	received := len(result.ResultSet.Results)
	if received == 0 {
		return nil
	}
	*courses = append(*courses, result.ResultSet.Results...)
	if received == result.FiltersApplied.Limit {
		filter.Offset += result.FiltersApplied.Limit
		return fetchCourses(baseURL, filter, courses)
	}
	return nil
}
