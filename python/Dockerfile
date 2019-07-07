FROM python:2.7.16-slim-stretch

COPY ./python/dingdian.py /python/

COPY ./python/LogUtils.py /python/

COPY ./python/requirements.txt /python/

WORKDIR /python

RUN pip install -r requirements.txt

MAINTAINER amanohiakru "1060996790@qq.com"

CMD ["python", "/python/dingdian.py"]