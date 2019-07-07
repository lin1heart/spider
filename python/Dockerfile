FROM python:2.7.16-slim-stretch

COPY ./spider/python/dingdian.py /python/

COPY ./spider/python/LogUtils.py /python/

COPY ./spider/python/requirements.txt /python/

WORKDIR /python

RUN pip install -r requirements.txt

MAINTAINER amanohiakru "1060996790@qq.com"

CMD ["python", "/python/dingdian.py"]