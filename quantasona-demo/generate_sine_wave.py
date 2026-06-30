import wave
import struct
import math

sample_rate = 16000.0
duration_seconds = 5.0
frequency = 440.0 # A4 note to simulate human voice frequency band

with wave.open('real_voice_sample.wav', 'w') as obj:
    obj.setnchannels(1) # mono
    obj.setsampwidth(2)
    obj.setframerate(sample_rate)
    
    for i in range(int(sample_rate * duration_seconds)):
        value = int(32767.0*math.cos(frequency*math.pi*float(i)/float(sample_rate)))
        data = struct.pack('<h', value)
        obj.writeframesraw(data)
