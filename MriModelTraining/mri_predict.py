from typing import Tuple
from MLPtrainer import Network

MODEL_DIR = "./mri_prediction_model/mri_model.json"
neural_network = Network.loadModel(MODEL_DIR)

def mriPredict(values) -> int:
    # (educ, mmse, brain_volume)
    out = neural_network.forwardpass(values)
    return int(round(out[-1][1][0]))