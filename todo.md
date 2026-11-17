### ARCHETECTURe

I say -> "hello"

send audio as a stream into Go
Go passes through audio into a speech detector
If there is speech, it will call the endpoint for stt
python runs the correct model and returns the text.
if the optimize speech is selected after the raw text is returned from the transcriber
it will pass the text into an llm with a prompt saying, around 90% of the words should be correct, go through each sentance and Identify the meaning of a sentance, if there is a word that is incorrect find the most likely word that was meant to be said, apply a confidence value aswell.
The prompt gets sent back to go, and go identifies the correct context.
The prompt with context is then sent to the llm in python.
For tooling, show 3-5 common tools, ex. web search. and a specific tooling tool
the llm with respond in a stream. While it streams the stream, return the text to go.
In Go, identify any parts that should not be spoken. Ex. a script or an attached file.
Then pass the stream back into python tts endpoint.
Python turns the text stream into an audio stream
Go plays the audio and shows whats being read.


