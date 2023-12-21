import base64 
import sys
import traceback
import copy
from PIL import Image

def read_image(image_file_path, base64_format=True):
    with open(image_file_path, "rb") as file:
        data = file.read()
        if not base64_format:
            return data

        base64_data = base64.b64encode(data)
        return base64_data 

def write_image(image_data, image_file_path, base64_format=True):
    if base64_format:
        data = base64.b64decode(image_data)
    else:
        data = image_data
    f = io.BytesIO(data)
    Image.open(f)
    os.makedirs(os.path.dirname(image_file), exist_ok=True)
    image.save(image_file_path)


class CVSDK:
    def __init__(self, gpu_id: int = -1, **kwargs):
        self.gpu_id = -1
        self.init_kwargs = {}
        self.sdk = None

    def init(self, model_dir: str, gpu_id: int = -1, **kwargs):
        try:
            if not self.sdk:
               del self.sdk
            self.sdk = None
            self.gpu_id = gpu_id
            self.init_kwargs = copy.copy(kwargs)

            if "sdk_dir" in kwargs:
                sdk_dir = kwargs.get("sdk_dir")
                del kwargs["sdk_dir"]
                if sdk_dir not in sys.path:
                    sys.path.append(sdk_dir)

            # init sdk isinstance
            # self.sdk = None

        except ValueError as err:
            return {
                    "err" : f"CVSDK::__init__: Unexpected ValueError {err=}, {type(err)=}, {traceback.format_exc()}",
                    }
        except Exception as err:
            return {
                    "err" : f"CVSDK::__init__: Unexpected Exception {err=}, {type(err)=}, {traceback.format_exc()}",
                    }
        except:
            err = sys.exc_info()[0]
            return {
                    "err" : f"CVSDK::__init__: Unexpected Exception {err=}, {type(err)=}, {traceback.format_exc()}",
                    }
        return {}

    
    def do(self, **kwargs):
        try:
            arg1 = kwargs.get("arg1", "")
            arg2 = kwargs.get("arg2", "")
            #todo sdk call
            return {
                    "result": "ok"
                    }
        except ValueError as err:
            return {
                    "err" : f"CVSDK::Do: Unexpected ValueError {err=}, {type(err)=}, {traceback.format_exc()}",
                    }
        except Exception as err:
            return {
                    "err" : f"CVSDK::Do: Unexpected Exception {err=}, {type(err)=}, {traceback.format_exc()}",
                    }
        except:
            err = sys.exc_info()[0]
            return {
                    "err" : f"CVSDK::Do: Unexpected Exception {err=}, {type(err)=}, {traceback.format_exc()}",
                    }

if __name__ == "__main__":
    print('sys.path %s' % sys.path)
    sdk = CVSDK()
    sdk.init("")
    resp = sdk.do()
    if 'err' in resp:
        err = resp['err']
        print(f'do err: {err}')
    print(f'{resp["result"]}')
