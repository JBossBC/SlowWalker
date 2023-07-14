import csv
import math
import os.path
import re

from lxml import etree


class outputParams:
    def __init__(self, targetFile, targetHeader):
        self.targetFile = targetFile
        if os.path.exists(targetFile):
            raise "输出文件已经存在"
        pointIndex = -1
        for index in reversed(range(len(str(targetFile)))):
            if targetFile[index] == '.':
                pointIndex = index
                break
        if pointIndex == -1:
            raise "cant find the file type"
        self.fileName = targetFile[:pointIndex]
        self.fileType = targetFile[pointIndex + 1:]
        self.targetHeader = targetHeader


class inputParams:
    # the fileterFiles should be dict
    def __init__(self, file, filterFiles):
        if not os.path.exists(file):
            raise "未找到该文件"
        self.file = file
        if not isinstance(filterFiles, set) and filterFiles is not None:
            raise "the filterFiles should be dict"
        self.filterFiles = filterFiles
        if os.path.isdir(file):
            self.isDir = True
        else:
            self.isDir = False

    def isFilterFile(self, file):
        if self.filterFiles is None:
            return False
        for out in range(self.filterFiles):
            if str(out) == str(file):
                return True

        return False


'''
  handle:
     tablePath: 一行数据的开始元素
     unitSelect: 一列的数据选择(如果传入字符串类型，则判定为每一行中每一列数据都用相同的规则)
     pageBegin: 一个文件中真实数据从多少个元素开始
     pageNumber: 一个文件中存在多少个数据行
     pageOffset: 数据行的偏移量
'''
MAX_CSVROWS = 1048576

DEFAULT_EXTRAFILENAME = "2"


def isNumber(s):
    pattern = r'^[-+]?(\d+(\.\d*)?|\.\d+)([eE][-+]?\d+)?$'
    return bool(re.match(pattern, s))


def _sortByNumber(item):
    return int(item)


class handle:
    data = []
    numberSort = False

    def __init__(self, tableXPath, unitSelect, pageNumber, pageBegin, pageOffset):
        if isinstance(unitSelect, str):
            self.equalUnit = True
        else:
            self.equalUnit = False
        if pageBegin < 0:
            raise "pageBegin cant less than zero"
        self.unitSelect = unitSelect
        self.tableXPath = tableXPath
        self.pageNumber = pageNumber
        self.pageBegin = pageBegin
        self.pageOffset = pageOffset

    def wash(self, inputData, outputData):
        if not isinstance(inputData, inputParams) or not isinstance(outputData, outputParams):
            raise "input params cant correspond,should be input and output type"

        if not inputData.isDir:
            self.handleFile(inputData.file)
        ## cant support the Loop traversal
        try:
            for root, dirs, files in os.walk(inputData.file):
                if not self.numberSort and files is not None and isNumber(files[0]):
                    self.numberSort = True
                files = files if not self.numberSort else sorted(files, key=_sortByNumber)
                for file in files:
                    if inputData.filterFiles is not None and str(file) in inputData.filterFiles:
                        continue
                    self._handleFile(str(inputData.file + "\\" + file))
        except Exception as e:
            print(str("运行到{}文件时出现异常:{}".format(file, e.__str__())) + "\n")

        self._finalizerData(outputData)

    def _finalizerData(self, outputData):
        if self.data is None:
            return
        with open(outputData.targetFile, "w+", encoding="utf-8", newline="") as f:
            w = csv.writer(f)
            rowsNumber = 1
            w.writerow(outputData.targetHeader)
            for index in range(len(self.data)):
                w.writerow(self.data[index])
                rowsNumber = rowsNumber + 1
                if rowsNumber >= MAX_CSVROWS:
                    if not f.closed:
                        f.close()
                    f = open(
                        str(outputData.fileName) + "({}).".format(DEFAULT_EXTRAFILENAME) + str(outputData.fileType),
                        encoding="utf-8", newline="")
                    w = csv.writer(f)
                    w.writerow(outputData.targetHeader)
                    rowsNumber = 1
                    DEFAULT_EXTRAFILENAME = DEFAULT_EXTRAFILENAME + 1

    def _handleFile(self, file):
        with open(r'{}'.format(file), 'r', encoding="utf-8") as f:
            res = f.read()
        if os.stat(file).st_size == 0:
            print("{} 文件为空,请检查\n".format(file))
            return
        html = etree.HTML(res)
        lineLength = len(html.xpath(self.tableXPath))
        if lineLength >= self.pageNumber + self.pageOffset * self.pageNumber:
            lineLength = self.pageNumber
        elif lineLength < self.pageNumber + self.pageOffset * self.pageNumber:
            lineLength = int(math.floor(float(lineLength / (self.pageOffset + 1))))
        times = 0

        if self.equalUnit:
            times = 1
        else:
            times = len(self.unitSelect)
        for row in range(self.pageBegin, self.pageBegin + (self.pageOffset + 1) * lineLength, self.pageOffset + 1):
            result = []
            for i in range(times):
                try:
                    if self.unitSelect == "":
                        continue
                    infos = html.xpath(self.tableXPath + "[{}]".format(row) + "/" + self.unitSelect[i])
                    info = ""
                    for information in infos:
                        if isinstance(information, str):
                            info = info + " " + information
                    # for information in infos:
                    # info.join()
                    if info is None or info == "":
                        info = "-"
                    result.append(info)
                except IndexError:
                    info = "-"
                    result.append(info)
                    pass
                except Exception as e:
                    print("解析文件的时候出现异常,xpath解析后的数据:",infos,"  使用的xpath规则:",self.tableXPath + "[{}]".format(row) + "/" + self.unitSelect[i])
                    raise "程序异常:{}".format(e.__str__())

            self.data.append(result)


if __name__ == "__main__":
    # handler = handle(r'//tr[@class="m_cen"]',
    #                  ["td[1]/a[1]/text()", "td[2]/text()", "td[3]/text()", "td[4]/text()", "td[5]/a[1]/text()", "td[6]/text()",
    #                   "td[7]/a[1]/text()", "td[8]/text()", "td[9]/text()"], 20, 0, 0)
    # handler.wash(inputData=inputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\数据\会员账号2023.3.20-2023.3.29',set("0")),outputData=outputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\数据\会员账号2023.3.20-2023.3.29.csv',["状态","真实姓名","登入账号","层级","系统额度","盘口","所属代理商","新增日期","最后登录日期"]))
    # handler = handle(r'//tr[@class="m_cen"]',
    #                  ["td[1]/text()", "td[2]/text()", "td[3]/text()", "td[4]/text()", "td[5]/node()",
    #                   "td[6]/span[1]/text()",
    #                   "td[7]/node()", "td[8]/font[1]/text()", "td[9]/text()", "td[10]/text()", "td[11]/text()"], 30, 0,
    #                  0)
    # handler.wash(
    #     inputData=inputParams(r'D:\work\湖南\5972.com\数据\入款记录(2023.2.1-2023.4.16)(全部状态)', set("0")),
    #     outputData=outputParams(r'D:\work\湖南\5972.com\数据\入款记录(2023.2.1-2023.4.16)(全部状态).csv',
    #                             ["层级", "订单号", "代理商", "会员账号", "会员银行账号", "存入金额", "存入银行账号", "状态", "是否首存", "操作者", "时间"]))
    # handler = handle(r'//tr[@class="m_cen"]',
    #                  ["td[1]/text()", "td[2]/text()", "td[3]/text()", "td[4]/text()", "td[5]/a[1]/text()"], 579, 0,
    #                  0)
    # handler.wash(
    #     inputData=inputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\数据\有效会员记录2023.3.20-2023.3.29', set("0")),
    #     outputData=outputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\数据\有效会员记录2023.3.20-2023.3.29.csv',
    #                             ["排序", "账号", "登录账号", "名称", "人数"]))
    # handler = handle(r'//tr[@class="m_cen"]',
    #                  ["td[1]/text()", "td[2]/text()", "td[3]/text()", "td[4]/text()", "td[5]/text()"], 100, 0,
    #                  0)
    # handler.wash(
    #     inputData=inputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\数据\自动稽查2023.3.20-2023.3.31', set("0")),
    #     outputData=outputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\数据\自动稽查2023.3.20-2023.3.31.csv',
    #                             ["层级", "账号", "讯息", "登入时间", "IP位置"]))
    # handler = handle(r'//tr[@class="m_cen"]',
    #                  [ "td[2]/text()", "td[3]/text()", "td[4]/text()", "td[5]/text()", "td[6]/text()",
    #                   "td[7]/text()",
    #                   "td[8]/text()", "td[9]/text()"], 15, 0,
    #                  0)
    # handler.wash(
    #     inputData=inputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\数据\邀请好友记录2023.3.20-2023.3.29', set("0")),
    #     outputData=outputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\数据\邀请好友记录2023.3.20-2023.3.29.csv',
    #                             ["会员账号", "会员姓名", "好友账号", "好友姓名", "注册时间","注册IP","最终登入时期","參與優惠"]))
    # handler = handle(r'//tr[@class="m_cen"]',
    #                  [ "td[1]/text()","td[2]/text()", "td[3]/text()", "td[4]/text()", "td[5]/text()", "td[6]/text()"], 50, 0,
    #                  0)
    # handler.wash(
    #     inputData=inputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\补数据\登录日志(2023.3.18-2023.4.1)', set("0")),
    #     outputData=outputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\补数据\登录日志(2023.3.18-2023.4.1).csv',
    #                             ["层级", "账号", "讯息", "登入时间","IP位置","设备"]))
    # handler = handle(r'//tr[@class="m_cen"]',
    #                  ["td[2]/text()", "td[3]/text()", "td[4]/text()", "td[5]/text()", "td[8]/text()", "td[11]/text()",
    #                   "td[14]/text()","td[15]/text()","td[16]/text()","td[17]/text()","td[18]/text()"],
    #                  100, 0,
    #                  0)
    # handler.wash(
    #     inputData=inputParams(r'D:\work\湖南\5972.com\数据\会员等级信息', set("0")),
    #     outputData=outputParams(r'D:\work\湖南\5972.com\数据\会员等级信息.csv',
    #                             ["会员账号", "注册日期", "最近登录日期", "会员等级", "真人等级", "电子等级","总存款","总投注金额",""
    #                                                                                             "視訊投注金額","电子投注金额","數據更新時間"]))
    # handler = handle(r'//table[@class="table table-hover table-bordered"]/tr',
    #                  ["td[2]/text()","td[4]/a[1]/text()", "td[5]/text()", "td[6]/text()", "td[7]/text()", "td[8]/font[1]/text()","td[9]/font[1]/text()",
    #                   "td[10]/font[1]/text()","td[11]/font[1]/text()","td[12]/text()","td[13]/text()","td[14]/text()","td[15]/text()",
    #                   "td[16]/text()","td[17]/text()","td[18]/text()","td[19]/a[1]/text()","td[20]/text()","td[21]/text()","td[22]/font[1]/text()","td[23]/b[1]/text()"], 20, 1,
    #                  0)
    # handler.wash(
    #     inputData=inputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\补数据\登录日志(2023.3.18-2023.4.1)', set("0")),
    #     outputData=outputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\补数据\登录日志(2023.3.18-2023.4.1).csv',
    #                             ["ID", "用户名称", "昵称", "LV","在线","测试号","异常","冻结","昵称修改禁用","stream","上级","账户余额","可用资金",
    #                              "冻结资金","可用余额","BB余额","BB可换","注册日期","最后登录","生日","资金锁","转移"]))

    # handler = handle(r'//tr[@class="m_cen"]',
    #                  ["td[1]/text()", "td[2]/a[1]/text()", "td[3]/button[1]/text()", "td[4]/text()",
    #                   "td[5]/a[1]/text()", "td[6]/button[1]/text()"], 11, 0,
    #                  0)
    # handler.wash(
    #     inputData=inputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\数据\出入账目汇总记录2023.3.20-2023.3.29', set("0")),
    #     outputData=outputParams(r'C:\Users\Xiyang\Desktop\work\湖南\桃源\5972\数据\出入账目汇总记录2023.3.20-2023.3.29.csv',
    #                             ["收入", "收入明细", "", "支出", "支出明细", ""]))

    # handler = handle(r'//tr[@class="m_cen"]',
    #                  ["td[1]/a[1]/text()", "td[2]/font[1]/text()", "td[3]/text()", "td[4]/a[1]/text()",
    #                   "td[5]/text()", "td[6]/a[1]/text()",
    #                   "td[7]/font[2]/text()","td[8]/text()","td[9]/text()",
    #                   "td[10]/text()","td[11]/text()","td[12]/text()",
    #                   "td[13]/text()","td[14]/text()","td[15]/text()"], 100, 0,
    #                  0)
    # handler.wash(
    #     inputData=inputParams(r'D:\work\湖南\5972.com\数据\代理商账号列表', set("0")),
    #     outputData=outputParams(r'D:\work\湖南\5972.com\数据\代理商账号列表.csv',
    #                             ["状态", "类型", "代理名称", "代理账号", "登入账号", "所属总代理","会员数","代理占成","总代理占成",
    #                              "新增日期","状况"]))
    # handler = handle(r'//tr[@class="m_cen"]',
    #                  ["td[1]/text()", "td[2]/text()", "td[3]/text()", "td[4]/text()",
    #                   "td[5]/text()", "td[6]/text()",
    #                   "td[7]/text()","td[8]/text()","td[9]/text()",
    #                   "td[10]/text()","td[11]/text()","td[12]/text()",
    #                   "td[13]/text()","td[14]/text()","td[15]/text()","td[16]/text()","td[17]/text()",
    #                   "td[18]/text()","td[19]/text()","td[20]/text()","td[21]/text()","td[22]/a[1]/text()",
    #                   "td[23]/text()","td[24]/text()","td[25]/text()","td[26]/text()","td[27]/text()","td[28]/text()","td[29]/a[1]/text()"], 100, 0,
    #                  0)
    # handler.wash(
    #     inputData=inputParams(r'D:\work\湖南\5972.com\数据\代理退佣(2022.12)', set("0")),
    #     outputData=outputParams(r'D:\work\湖南\5972.com\数据\代理退佣(2022.12).csv',
    #                             ["代理账号", "名称", "有效会员", "派彩", "当期", "体育退佣比例","沙巴退佣比例","彩票退佣比例","视讯退佣比例",
    #                              "电子退佣比例","已获退佣","状态","动作","出款银行资料","备注"]))
    handler = handle(r'//tr[@class="m_cen"]',
                     ["td[1]/a[1]/text()","td[2]/text()", "td[3]/text()", "td[4]/text()", "td[5]/a[1]/text()", "td[6]/text()", "td[7]/a[1]/text()",
                      "td[8]/text()", "td[9]/text()", "td[10]/text()"],
                     50, 0,
                     0)
    handler.wash(
        inputData=inputParams(r'D:\work\湖南\5972.com\数据\会员列表', set("0")),
        outputData=outputParams(r'D:\work\湖南\5972.com\数据\会员列表.csv',
                                ["状态", "真实姓名", "登入账号", "层级", "系统额度", "盘口", "所属代理商", "新增日期", "最后登录日期","状况"]))
