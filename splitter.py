import acoustid
import chromaprint
from fuzzywuzzy import fuzz
import pyaudio
import asyncio

print("Choose your Audio Input:\n")
p = pyaudio.PyAudio()
info = p.get_host_api_info_by_index(0)
numdevices = info.get('deviceCount')
for i in range(0, numdevices):
    if (p.get_device_info_by_host_api_device_index(0,i).get('maxInputChannels')) > 0:
            print("Input Device Id: ", i , " - ", p.get_device_info_by_host_api_device_index(0,i).get('name'))
audio_input = input('Device ID---> ')
p.input_device_index=audio_input
CHUNK = 1024
WIDTH = 2
CHANNELS = 2
RATE = 44100
RECORD_SECONDS = 0.1

duration, fp_encoded = acoustid.fingerprint_file('lap.mp3')
fingerprint, version = chromaprint.decode_fingerprint(fp_encoded)

async def readMic():
    stream = p.open(format=p.get_format_from_width(WIDTH),
            channels=CHANNELS,
            rate=RATE,
            input=True,
            output=True,
            frames_per_buffer=CHUNK)
    frames = []
    for i in range(0, int(RATE/CHUNK * RECORD_SECONDS)):
        data = stream.read(CHUNK)
        frames.append(data)
    stream.stop_stream()
    stream.close()
    p.terminate()
    wf = wave.open("temp.wav", 'wb')
    wf.setnchannels(CHANNELS)
    wf.setsampwidth(p.get_sample_size(FORMAT))
    wf.setframerate(RATE)
    wf.writeframes(b''.join(frames))
    wf.close()
    p.close

async def compareFingerprints():
    sample_duration, sample_fp_encoded = acoustid.fingerprint_file('temp.wav')
    sample_fingerprint, sample_version = chromaprint.decode_fingerprint(sample_fp_encoded)
    similarity = fuzz.ratio(sample_fingerprint, fingerprint)

async def autoKeyboard():
    asyncio.run(readMic())
    asyncio.run(compareFingerprints())
