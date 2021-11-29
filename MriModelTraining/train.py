import os
from tensorflow.keras.preprocessing.image import ImageDataGenerator
import tensorflow as tf
import matplotlib.pyplot as plt
import numpy as np
from tensorflow.python import keras
from tensorflow.python.util.decorator_utils import validate_callable
import cv2

folder_training_data = "./Alzheimer_s Dataset/train"
folder_testing_data = "./Alzheimer_s Dataset/test"
img_size = (200,200)
bach_size = 52
epochs= 18
learning_rate = 0.001

trainer = ImageDataGenerator(rescale=1/255)
tester = ImageDataGenerator(rescale=1/255)

trainer_data = trainer.flow_from_directory(
    folder_training_data,
    target_size=img_size,
    batch_size=bach_size,
    class_mode="binary"
)

tester_data = trainer.flow_from_directory(
    folder_testing_data,
    target_size=img_size,
    batch_size=bach_size,
    class_mode="binary"
)

print(trainer_data.class_indices)

model = tf.keras.models.Sequential([
    tf.keras.layers.Conv2D(16, (3,3), activation='relu', input_shape=(200,200,3)),
    tf.keras.layers.MaxPooling2D(2,2),
    tf.keras.layers.Conv2D(32, (3,3), activation='relu'),
    tf.keras.layers.MaxPooling2D(2,2),
    tf.keras.layers.Conv2D(64, (3,3), activation='relu'),
    tf.keras.layers.MaxPooling2D(2,2),
    tf.keras.layers.Flatten(),
    tf.keras.layers.Dense(512, activation='relu'),
    tf.keras.layers.Dense(1, activation='sigmoid')
])

model.compile(loss='binary_crossentropy', optimizer=tf.keras.optimizers.RMSprop(learning_rate=learning_rate), metrics=['accuracy'])
model.fit(
    trainer_data,
    steps_per_epoch=50,
    epochs=epochs,
    validation_data=tester_data
)
#print convergence
print(model.history.history['accuracy'])

target_dir = "ex_model"
if not os.path.exists(target_dir):
    os.mkdir(target_dir)
model.save(os.path.join(target_dir, "modelo.h5"))
model.save_weights(os.path.join(target_dir, "pesos.h5"))

