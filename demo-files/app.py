from flask import Flask, render_template, request
from flask_socketio import SocketIO, emit
import uvicorn
import socketio
import io
import json
import numpy as np
# import soundfile
import wave
import torchaudio
import torch
import time
from collections import defaultdict
from pathlib import Path
import pydub.effects as effects
from pydub import AudioSegment
from pydub.effects import normalize, compress_dynamic_range, low_pass_filter, high_pass_filter

from flask_cors import CORS
from seamless_communication.inference import Translator
from seamless_communication.streaming.dataloaders.s2tt import SileroVADSilenceRemover
import wave

### Init Config
MODEL_NAME = "seamlessM4T_v2_large"

config = {
  "device": "cuda:0",
  "model_name": MODEL_NAME,
  "vocoder_name": "vocoder_v2" if MODEL_NAME == "seamlessM4T_v2_large" else "vocoder_36langs",
  "target_lang": "cmn"
}

### Init Model
try:
    translator = Translator(
        config['model_name'],
        config['vocoder_name'],
        device=torch.device(config['device']),
        dtype=torch.float16,
    )

    print(f"Translator and model loaded to Device: {config['device']}")

except Exception as e:
    print(f"Error occuries when creating Translator: {e}")


### Init FastAPI App
app = Flask(__name__)
CORS(app)

### Init Async Server
sio = socketio.AsyncServer(
    async_mode="asgi",
    cors_allowed_origins="*",
)

socket_app = socketio.ASGIApp(
    sio, 
    app,
    static_files={
    '/': 'templates/index.html',
    }
)

@sio.on('translate')
async def translate_speech(sid, message):
    print(message)
    speech_local_file_path = message['data']['speech_local_file_path']
    target_lang = message['data']['target_lang']
    _start_time = time.time()

    print(f'Target Language is {target_lang}')

    text_output, speech_output = translator.predict(
      input=speech_local_file_path, # can be file or AudioTensor
      task_str="s2st",
      tgt_lang=target_lang,
    )

    print(f"Translated text in {target_lang}: {text_output[0]}")
    print(f"Time collasped for {target_lang}: {time.time()-_start_time}")

    await sio.emit('display_speech_transcription', {'data': str(text_output[0])})

background_task_started = False

async def background_task():
    """Example of how to send server generated events to clients."""
    # count = 0
    while True:
        await sio.sleep(10)


@sio.on('my_event')
async def test_message(sid, message):
    await sio.emit('my_response', {'data': message['data']}, room=sid)


@sio.on('my_broadcast_event')
async def test_broadcast_message(sid, message):
    await sio.emit('my_response', {'data': message['data']})


@sio.on('join')
async def join(sid, message):
    await sio.enter_room(sid, message['room'])
    await sio.emit('my_response', {'data': 'Entered room: ' + message['room']},
                   room=sid)


@sio.on('leave')
async def leave(sid, message):
    await sio.leave_room(sid, message['room'])
    await sio.emit('my_response', {'data': 'Left room: ' + message['room']},
                   room=sid)


@sio.on('close room')
async def close(sid, message):
    await sio.emit('my_response',
                   {'data': 'Room ' + message['room'] + ' is closing.'},
                   room=message['room'])
    await sio.close_room(message['room'])


@sio.on('my_room_event')
async def send_room_message(sid, message):
    await sio.emit('my_response', {'data': message['data']},
                   room=message['room'])

@sio.on('connect')
async def connect(sid, environ):
    global background_task_started
    if not background_task_started:
        sio.start_background_task(background_task)
        background_task_started = True
    await sio.emit('my_response', {'data': 'Connected', 'count': 0}, room=sid)


@sio.on('disconnect')
async def disconnect(sid):
    print('Client disconnected')
    await sio.disconnect(sid)


@sio.on('audio_data')
async def handle_audio_data(sid, data):
    # Process the received audio data asynchronously
    print('Received audio data:', len(data))

    wav_file_bytesIO = io.BytesIO(data)
    data, sample_rate = torchaudio.load(wav_file_bytesIO)
    # # Here can try to do the translator.predict, with passing "data" as "input" (because input supports AudioTensor)
    print(f"Sample Rate is {sample_rate}")

    target_lang = 'eng' #message['data']['target_lang']
    _start_time = time.time()

    print(f'Target Language is {target_lang}')

    text_output, speech_output = translator.predict(
      input=data, # can be file or AudioTensor
      task_str="S2ST", # s2st
      tgt_lang=target_lang,
      sample_rate=48000, # suggest to pass sample rate for better performance
    )

    print("Save audio file for verification")
    audio_tensor = speech_output.audio_wavs[0][0].to(torch.float32).cpu()

    translated_audio_path = '/home/ubuntu/seamless-streaming/seamless_server/audio_output/test.wav'

    torchaudio.save(
        translated_audio_path, 
        audio_tensor, 
        # speech_output.sample_rate
        16000 # 16000 works
        )

    print(f"Translated text in {target_lang}: {text_output[0]}")
    print(f"Time collasped for {target_lang}: {time.time()-_start_time}")

    await sio.emit('display_speech_transcription', {'data': str(text_output[0])})

    # # Refer to https://github.com/amsehili/auditok/issues/47#issuecomment-1813975148
    # # translated_bytes_buffer = (audio_tensor.numpy() * 32767).astype(np.int16).tobytes()
    # audio_numpy = audio_tensor.numpy()

    _start_time = time.time()
    audio_buffer = io.BytesIO()

    """
        Convert Audio Tensor to Bytes and send back to Frontend.
    """
    """
    Wave works, but audio quality is too bad - NOT RECOMMENDED
    """
    # wave_writer = wave.open(audio_buffer, 'wb')
    # wave_writer.setnchannels(1)  # Assuming mono audio
    # wave_writer.setsampwidth(4)  # 32-bit float
    # wave_writer.setframerate(16000)  # Sample rate (adjust as needed)
    # wave_writer.writeframes(audio_numpy.astype(np.float32).tobytes())
    # wave_writer.close()

    """
       Try to use AudioSegment
    """
    # Create an AudioSegment from the numpy array
    # audio_segment = AudioSegment(
    #     data=audio_numpy.tobytes(),
    #     sample_width=4,  # 32-bit float
    #     frame_rate=16000,  # Sample rate (adjust as needed) default 16000
    #     channels=1  # Assuming mono audio, 1 for mono, 2 for stereo
    # )

    audio_segment = AudioSegment.from_file(translated_audio_path, 'wav')
    audio_segment = audio_segment.set_frame_rate(16000)
    audio_segment = audio_segment.set_channels(1)
    audio_segment = audio_segment.set_sample_width(4)

    # 归一化 (normalize)
    # 动态范围压缩 (compress_dynamic_range)
    # 低通滤波器 (low_pass_filter)
    # 高通滤波器 (high_pass_filter)
    # 噪声门限滤波 (strip_silence)
    # 音量调整 (apply_gain)

    # Apply audio processing techniques
    # processed_audio = normalize(audio_segment, headroom=-3.0)  # Normalize to -3/-1 dBFS
    # processed_audio = compress_dynamic_range(processed_audio, threshold=-30.0, ratio=8.0)  # Compress dynamic range
    # processed_audio = low_pass_filter(processed_audio, 5000)  # Apply low-pass filter at 18 kHz
    # processed_audio = high_pass_filter(processed_audio, 80)  # Apply high-pass filter at 80 Hz
    # processed_audio = effects.strip_silence(processed_audio, silence_thresh=-30)

    # # Decrease Volume
    # # processed_audio = processed_audio.apply_gain(-3)  # Decrease by 6 dB

    # # Export the processed AudioSegment to a BytesIO buffer in WAV format
    # audio_buffer = io.BytesIO()
    # processed_audio.export(audio_buffer, format='wav')

    # Export the AudioSegment to a BytesIO buffer in WAV format
    audio_segment.export(audio_buffer, format='wav')

    # # Get the audio data as bytes
    audio_bytes = audio_buffer.getvalue()
    
    print(f"Time collasped for Audio propagation: {time.time()-_start_time}")

    await sio.emit('play_translated_audio', audio_bytes) 


if __name__ == '__main__':
    uvicorn.run(socket_app, host='127.0.0.1', port=5000)
