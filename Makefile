# run gomobile bind -target macos ./worker

clean:
	rm -rf Worker.xcframework

swift:
	gomobile bind -target macos ./worker