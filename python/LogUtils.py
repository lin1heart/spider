# coding=utf-8
import logging
import getpass
import sys


class Log(object):
    def __init__(self):
        self.user = getpass.getuser()
        self.logger = logging.getLogger(self.user)
        self.logger.setLevel(logging.DEBUG)

        # 日志文件名
        self.logFile = sys.argv[0][0:-3] + '.log'
        # self.formatter = logging.Formatter('%(asctime)-12s %(levelname)- 8s %(name)-10s %(message)-12s\r\n')
        self.formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s : %(message)s")
        # 日志显示到屏幕并输出到文档
        self.logHand = logging.FileHandler(self.logFile, encoding='utf-8')
        self.logHand.setFormatter(self.formatter)
        self.logHand.setLevel(logging.DEBUG)

        self.logHandSt = logging.StreamHandler()
        self.logHandSt.setFormatter(self.formatter)
        self.logHandSt.setLevel(logging.DEBUG)

        self.logger.addHandler(self.logHand)
        self.logger.addHandler(self.logHandSt)

    # 日志的5个级别对应5个函数
    def debug(self, msg):
        self.logger.debug(msg)

    def info(self, msg):
        self.logger.info(msg)

    def warn(self, msg):
        self.logger.warn(msg)

    def error(self, msg):
        self.logger.error(msg)

    def critical(self, msg):
        self.logger.critical(msg)
