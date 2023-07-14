import codecs
import csv
import json
import os


def AnalysisJSON(src, target):
    result =[]
    for root, dirs, files in os.walk(src):
        for file in files:
            data = ""
            with open(root+os.sep+file, "rb") as f:
                data = f.read()
                data = codecs.decode(data, 'utf-8-sig')
                print(data)
            jsonStr = data.split('\r\n\r\n',2)[1]
            if jsonStr.startswith('\ufeff'):
                jsonStr=jsonStr.removeprefix('\ufeff')
            jsonData = json.loads(jsonStr)["list"]
            result.extend(jsonData)
    with open(target,"w+",encoding="utf8",newline=''  ) as f:
        writer=csv.writer(f)
        rowHeader = []
        header=False
        for i in result:
            rowData = []
            for key in i:
                if not header:
                    rowHeader.append(key)
                rowData.append(i[key])
            if not header:
                header=True
                writer.writerow(rowHeader)
            writer.writerow(rowData)


