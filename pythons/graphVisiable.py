import copy

import networkx as nx
import pandas
import pandas as pd
from lxml import etree

graph = nx.DiGraph()

preData = "上层账号"
nowData = "账号"


# you should define the preData and nowData before using the VisibleGraph function
def VisibleGraph(src, target):
    global root
    global preList
    global html
    global userInfo
    if not str(src).endswith(".xlsx"):
        raise "the input file must be xlsx"
    df = pd.read_excel(src)
    df = df.dropna(axis=1)
    # df = df.drop(columns= ['Name', 'dtype'])
    df.apply(columnsHandle, axis=1)
    # in_degrees = dict(graph.in_degree())
    # zero_indegree_nodes = [n for n in in_degrees if in_degrees[n] == 0]
    # print(userInfo)
    for node in preList:
        html = etree.Element("html")
        root = etree.SubElement(html, 'ul')
        element = etree.SubElement(root, 'li')
        element.text = str(node)
        if userInfo.get(node) is not None:
            element.text = "{} 具体信息如下:{}".format(str(node), formatSerise(userInfo.get(node)))
        # element.text = "{} 具体信息如下:{}".format(node, userInfo.get(node))
        recursionAdd(node, element)
        with open(target + "\\" + str(node) + ".html", 'w+', encoding='utf8') as f:
            f.write(etree.tostring(html, pretty_print=True, encoding='unicode'))


preList = []

userInfo = dict()


def formatSerise(serise=pandas.Series):
    result = ""
    for key, value in serise.iteritems():
        if key == "Name" or key== "dtype":
            continue
        result = result + key + ":" + value + "     "
    return result


def columnsHandle(row):
    global preData
    global nowData
    global graph
    global preList
    global userInfo
    pre = row[preData]
    now = row[nowData]
    # del row["Name"]
    # del row["dtype"]
    if not graph.has_node(pre):
        graph.add_node(pre)
        preList.append(pre)
    if not graph.has_node(now):
        graph.add_node(now)
        userInfo[now] = copy.copy(row)
    if not graph.has_edge(pre, now):
        graph.add_edge(pre, now)


def recursionAdd(rootNode, element):
    nodes = graph.neighbors(rootNode)
    flag = False
    lists = None
    for node in nodes:
        if not flag:
            lists = etree.SubElement(element, 'ul')
            flag = True
        subElement = etree.SubElement(lists, 'li')
        subElement.text = str(node)
        if userInfo.get(node) is not None:
            subElement.text = "{} 具体信息如下:{}".format(node, userInfo.get(node))
        recursionAdd(node, subElement)


if __name__ == "__main__":
    VisibleGraph(r'C:\Users\Xiyang\Desktop\work\湖南\临澧\2017\可视化\会员列表--清洗数据(5).xlsx',
                 r'C:\Users\Xiyang\Desktop\work\湖南\临澧\2017\可视化')
