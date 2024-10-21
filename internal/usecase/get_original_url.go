package usecase

import (
	"context"
	"time"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/malwarescanner"
	"urlShortenerService/internal/infrastructure/shorturl"
	"urlShortenerService/internal/infrastructure/statistics"

	"github.com/golang/glog"
)

// GetOriginalURLCmd represents the function signature of the command that retrieves an original URL given a slug
type GetOriginalURLCmd func(ctx context.Context, shortURL string) (string, error)

func withMalwareScan(f GetOriginalURLCmd, malwareScanner malwarescanner.Scanner) GetOriginalURLCmd {
	return func(ctx context.Context, slug string) (string, error) {
		url, err := f(ctx, slug)
		if err != nil {
			return url, err
		}

		// Scan the URL for malware
		malwareScanResult := make(chan malwarescanner.MalwareScanResult, 1)
		go malwareScanner.Scan(context.Background(), url, malwareScanResult)

		timeout := time.After(time.Second) // Timeout of 1 second
		select {
		case malwareScanRes := <-malwareScanResult:
			switch malwareScanRes {
			case malwarescanner.MalwareScanResultClear:
				return url, nil
			case malwarescanner.MalwareScanResultDetected:
				return "", malwarescanner.ErrMalswareURL
			default:
				// If malware scanner errored for something else than a malware, we log but ignore the error
				glog.Warningf("malware scanner errored for [%s]", url)
			}
		case <-timeout:
			glog.Warningf("malware scanner timed out for [%s]", url)
		}

		return url, nil
	}
}

// getOriginalURL retrieves an original URL given a slug
func getOriginalURL(slugValidatorCmd command.SlugValidatorCmd, shortURLStore shorturl.Store, statisticsStore statistics.Store) GetOriginalURLCmd {
	return func(ctx context.Context, slug string) (string, error) {
		// Ensure slug validity to avoid useless query to store
		err := slugValidatorCmd(slug)
		if err != nil {
			return "", err
		}

		// Retrieves URL
		urlMapping, err := shortURLStore.Get(ctx, slug)
		if err != nil {
			return "", err
		}

		// Update statistics
		go func(url string) {
			err := statisticsStore.SetURL(context.Background(), url, statistics.StatisticTypeAccessed)
			if err != nil {
				glog.Errorf("failed to set [%s] statistics for [%s]: %w", statistics.StatisticTypeAccessed, url, err)
			}
		}(urlMapping.OriginalURL)

		return urlMapping.OriginalURL, nil
	}
}

// GetOriginalURLWithMalwareScanCmdBuilder builds the command that will retrieves an original URL and scan it for malware
func GetOriginalURLWithMalwareScanCmdBuilder(slugValidatorCmd command.SlugValidatorCmd, malwareScanner malwarescanner.Scanner,
	shortURLStore shorturl.Store, statisticsStore statistics.Store) GetOriginalURLCmd {
	return withMalwareScan(
		getOriginalURL(slugValidatorCmd, shortURLStore, statisticsStore),
		malwareScanner)
}

// ForceGetOriginalURLCmdBuilder builds the command that will retrieves an original URL bypassing scan for malware
func ForceGetOriginalURLCmdBuilder(slugValidatorCmd command.SlugValidatorCmd, shortURLStore shorturl.Store, statisticsStore statistics.Store) GetOriginalURLCmd {
	return getOriginalURL(slugValidatorCmd, shortURLStore, statisticsStore)
}
