package twiutil

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/stretchr/testify/assert"
)

func readFileTweet(filename string) (*twitter.Tweet, error) {
	var tw twitter.Tweet
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return &tw, err
	}
	err = json.Unmarshal(data, &tw)
	if err != nil {
		return &tw, err
	}
	return &tw, nil
}

func createDummyTweet() (map[string]*twitter.Tweet, error) {
	fnames := []string{
		"tweet-gif.json",
		"tweet-linked.json",
		"tweet-movie.json",
		"tweet-photo.json",
		"tweet-quoted.json",
	}
	tweets := make(map[string]*twitter.Tweet, len(fnames))
	for _, fname := range fnames {
		tweet, err := readFileTweet("test/" + fname)
		if err != nil {
			return nil, err
		}
		name := strings.TrimSuffix(strings.TrimPrefix(fname, "tweet-"), ".json")
		tweets[name] = tweet
	}
	return tweets, nil
}

func TestHasMedia(t *testing.T) {
	tweets, err := createDummyTweet()
	if err != nil {
		assert.Fail(t, "cannot create dummy tweet.")
	}

	type args struct {
		tweet *twitter.Tweet
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "gif",
			args: args{tweets["gif"]},
			want: true,
		},
		{
			name: "linked",
			args: args{tweets["linked"]},
			want: false,
		},
		{
			name: "movie",
			args: args{tweets["movie"]},
			want: true,
		},
		{
			name: "photo",
			args: args{tweets["photo"]},
			want: true,
		},
		{
			name: "quoted",
			args: args{tweets["quoted"]},
			want: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, HasMedia(tt.args.tweet))
	}
}

func TestFindUrlAll(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "normal",
			args: args{url: "asdfhttps://twitter.com/FloodSocial/status/861627479294746624giuewr"},
			want: []string{"https://twitter.com/FloodSocial/status/861627479294746624"},
		},
		{
			name: "query",
			args: args{url: "asdfghhttps://twitter.com/i/statuses/861627479294746624?s=20"},
			want: []string{"https://twitter.com/i/statuses/861627479294746624"},
		},
		{
			name: "short",
			args: args{url: "asdfhhttp://twitter.com/status/861627479294746624?s=afw"},
			want: []string{"http://twitter.com/status/861627479294746624"},
		},
		{
			name: "multi",
			args: args{url: "https://twitter.com/FloodSocial/status/1440105622737854464?s=20asdfhhttp://twitter.com/status/861627479294746624?s=afw"},
			want: []string{"https://twitter.com/FloodSocial/status/1440105622737854464", "http://twitter.com/status/861627479294746624"},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, FindUrlAll(tt.args.url))
	}
}

func TestFindIdAll(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			name: "normal",
			args: args{url: "asdfhttps://twitter.com/FloodSocial/status/861627479294746624giuewr"},
			want: []int64{861627479294746624},
		},
		{
			name: "query",
			args: args{url: "asdfghhttps://twitter.com/i/statuses/861627479294746624?s=20"},
			want: []int64{861627479294746624},
		},
		{
			name: "short",
			args: args{url: "http://twitter.com/status/861627479294746624"},
			want: []int64{861627479294746624},
		},
		{
			name: "multi",
			args: args{url: "https://twitter.com/FloodSocial/status/1440105622737854464?s=20\nasdfhhttp://twitter.com/status/861627479294746624?s=afw"},
			want: []int64{1440105622737854464, 861627479294746624},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, FindIdAll(tt.args.url))
	}
}

func TestFindId(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "normal",
			args: args{url: "asdfhttps://twitter.com/FloodSocial/status/861627479294746624giuewr"},
			want: 861627479294746624,
		},
		{
			name: "query",
			args: args{url: "asdfghhttps://twitter.com/i/statuses/861627479294746624?s=20"},
			want: 861627479294746624,
		},
		{
			name: "short",
			args: args{url: "asdfhhttp://twitter.com/status/861627479294746624?s=afw"},
			want: 861627479294746624,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, FindId(tt.args.url))
	}
}

func TestGetMediaUrls(t *testing.T) {
	tweets, err := createDummyTweet()
	if err != nil {
		assert.Fail(t, "cannot create dummy tweet.")
	}

	type args struct {
		tweet *twitter.Tweet
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "gif",
			args: args{tweets["gif"]},
			want: []string{
				"https://video.twimg.com/tweet_video/DBMDLy_U0AAqUWP.mp4",
			},
		},
		{
			name: "linked",
			args: args{tweets["linked"]},
			want: []string{},
		},
		{
			name: "movie",
			args: args{tweets["movie"]},
			want: []string{
				"https://video.twimg.com/ext_tw_video/869317980307415040/pu/vid/720x1280/octt5pFbISkef8RB.mp4",
			},
		},
		{
			name: "photo",
			args: args{tweets["photo"]},
			want: []string{
				"https://pbs.twimg.com/media/C_UdnvPUwAE3Dnn.jpg",
				"https://pbs.twimg.com/media/C_UdnvPVYAAZbEs.jpg",
				"https://pbs.twimg.com/media/C_Udn2UUQAADZIS.jpg",
				"https://pbs.twimg.com/media/C_Udn4nUMAAgcIa.jpg",
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, GetMediaUrls(tt.args.tweet))
	}
}

func TestGetMediaTypes(t *testing.T) {
	tweets, err := createDummyTweet()
	if err != nil {
		assert.Fail(t, "cannot create dummy tweet.")
	}

	type args struct {
		tweet *twitter.Tweet
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "gif",
			args: args{tweets["gif"]},
			want: []string{"animated_gif"},
		},
		{
			name: "linked",
			args: args{tweets["linked"]},
			want: []string{},
		},
		{
			name: "movie",
			args: args{tweets["movie"]},
			want: []string{"video"},
		},
		{
			name: "photo",
			args: args{tweets["photo"]},
			want: []string{
				"photo",
				"photo",
				"photo",
				"photo",
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, GetMediaTypes(tt.args.tweet))
	}
}
