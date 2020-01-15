package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type StartTime struct {
	minutes int
	seconds int
	ms      int
}
type Track struct {
	number    int
	title     string
	artist    string
	startTime StartTime
}

type Album struct {
	name      string
	artist    string
	release   int
	genre     string
	tracks    *[]Track
	numTracks int
	file      string
}

const (
	ffmpegPath   = "ffmpeg"
	ffmpegOutput = true // show ffmpeg command output
	showArgs     = true // show ffmpeg arguments before executing
)

// Prints an Album in human-readable format
func (a *Album) print() {
	fmt.Printf("FILE: %s\n", a.file)
	fmt.Printf("TITLE: %s\n", a.name)
	fmt.Printf("ARTIST: %s\n", a.artist)
	fmt.Printf("GENRE: %s\n", a.genre)
	fmt.Printf("DATE: %d\n", a.release)
	fmt.Printf("TRACKS (%d):\n", a.numTracks)

	// print all tracks
	for _, track := range *a.tracks {
		fmt.Printf("    [%d] %s - %s [%02d:%02d:%02d]\n", track.number, track.title, track.artist, track.startTime.minutes, track.startTime.seconds, track.startTime.ms)
	}
}

// Converts a StartTime to a string
func (s *StartTime) toString() string {
	hours := s.minutes / 60
	minutes := s.minutes % 60

	output := fmt.Sprintf("%02d:%02d:%02d.%02d", hours, minutes, s.seconds, s.ms)
	return output
}

// Converts string of mm:ss:msms into a StartTime struct
func newTime(input string) StartTime {
	var time StartTime
	fmt.Sscanf(input, "%02d:%02d:%02d", &time.minutes, &time.seconds, &time.ms)

	return time
}

func escape(input string) string {
	return strings.ReplaceAll(input, "/", "-")
}

// Parses the given input file
func parseCue(input string) Album {

	thisAlbum := new(Album)
	var thisAlbumTracks []Track

	// open the cue file
	cueFd, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer cueFd.Close()

	// read the cue file...
	scanner := bufio.NewScanner(cueFd)
	for scanner.Scan() {
		// read the current line
		line := scanner.Text()

		// split the line at whitespace // TODO split at '""
		lineSplit := strings.Split(strings.TrimSpace(line), " ")

		/// TRACK PARSING ///
		switch lineSplit[0] {
		case "TRACK":
			var nextLine []string

			// extract the track number
			trackNum, _ := strconv.Atoi(lineSplit[1])
			thisTrack := Track{number: trackNum}

			// read the next line...
			scanner.Scan()
			nextLine = strings.Fields(scanner.Text())
			if nextLine[0] == "TITLE" {
				// everything after "TITLE" is the title...
				trackTitle := strings.Join(nextLine[1:len(nextLine)], " ")
				// remove '"'
				thisTrack.title = strings.Trim(trackTitle, "\"")
			}

			// read the next line...
			scanner.Scan()
			nextLine = strings.Fields(scanner.Text())
			if nextLine[0] == "PERFORMER" {
				trackArtist := strings.Join(nextLine[1:len(nextLine)], " ")
				thisTrack.artist = strings.Trim(trackArtist, "\"")
			}

			// read the next line...
			scanner.Scan()
			nextLine = strings.Fields(scanner.Text())
			if nextLine[0] == "INDEX" {
				thisTrack.startTime = newTime(nextLine[2])
			}

			// append these tracks to the album
			thisAlbumTracks = append(thisAlbumTracks, thisTrack)
			thisAlbum.tracks = &thisAlbumTracks
			thisAlbum.numTracks = len(thisAlbumTracks)
		case "REM":
			switch lineSplit[1] {
			case "GENRE":
				// genre is usually unquoted
				// so take everything after "GENRE":
				albumGenre := strings.Join(lineSplit[2:len(lineSplit)], " ")
				thisAlbum.genre = albumGenre
			case "DATE":
				var albumRelease int
				fmt.Sscanf(line, "REM DATE %4d", &albumRelease)
				thisAlbum.release = albumRelease
			}
		case "PERFORMER":
			var albumArtist string
			fmt.Sscanf(line, "PERFORMER %q", &albumArtist)
			thisAlbum.artist = albumArtist
		case "TITLE":
			var albumTitle string
			fmt.Sscanf(line, "TITLE %q", &albumTitle)
			thisAlbum.name = albumTitle
		case "FILE":
			var albumFile string
			fmt.Sscanf(line, "FILE %q WAVE", &albumFile)
			thisAlbum.file = strings.Replace(albumFile, "\"", "", -1)
		}

	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return *thisAlbum
}

func extractAlbum(album Album) {

	// for each track in the album...
	for index, track := range *album.tracks {
		var cmd *exec.Cmd

		outputName := fmt.Sprintf("%d. %s.flac", track.number, track.title)

		// print a message
		fmt.Printf("[%d] track %d > %q\n", index, track.number, outputName)

		// if the current track is the last one...
		if index == album.numTracks-1 {
			cmd = exec.Command(ffmpegPath, "-n",
				"-ss", track.startTime.toString(),
				"-i", album.file,
				escape(outputName))
		} else {
			// get the next track
			nextTrack := (*album.tracks)[index+1]
			// cut the file from the current track's start to the next track's start
			cmd = exec.Command(ffmpegPath, "-n",
				"-ss", track.startTime.toString(), "-i", album.file,
				"-to", nextTrack.startTime.toString(),
				escape(outputName))
		}

		if showArgs == true {
			fmt.Printf("%s\n", cmd.Args)
		}

		// run the ffmpeg command, capturing the output
		out, err := cmd.CombinedOutput()

		if ffmpegOutput == true {
			fmt.Printf("%s\n", out)
		}

		if err != nil {
			panic(err)
		}
	}
}

func main() {

	album := parseCue("./test.cue")
	album.print()
	print("extracting...\n")
	extractAlbum(album)
}
