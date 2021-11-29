import os
import numpy as np
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3' 
import tensorflow as tf
tf.compat.v1.logging.set_verbosity(tf.compat.v1.logging.ERROR) 
from keras.preprocessing.image import load_img, img_to_array
from tensorflow.python.keras.backend import set_session
from keras.models import load_model
from time import time

# tf.compat.v1.disable_eager_execution()


longitud, altura = 200, 200
target_dir = "./ex_model"
modelo = f'{target_dir}/modelo.h5'
pesos_modelo = f'{target_dir}/pesos.h5'
sess = tf.compat.v1.Session()
graph = tf.compat.v1.get_default_graph()
set_session(sess)
cnn = load_model(modelo)
cnn.load_weights(pesos_modelo)

def predict(file):
    assert os.path.exists(file)
    tf.compat.v1.global_variables_initializer()
    model = load_model(modelo)
    model.load_weights(pesos_modelo)
    file_data = load_img(file, target_size=(longitud, altura))
    file_data = img_to_array(file_data)
    print(file_data.shape)
    file_data = np.expand_dims(file_data, axis=0)
    result = model.predict(file_data)
    return result
    
os.system('cls')

if __name__ == "__main__":    
    print(predict("./Alzheimer_s Dataset/train/NonDemented/nonDem2.jpg"))
    # while True:
    #     os.system("clear")
    #     test_file = input("imagen: ")
    #     if test_file == "q":
    #         break
    #     else:
    #         start = time()
    #         print(predict(test_file))
    #         print(f"Tiempo: {time() - start}")
    #         input()
