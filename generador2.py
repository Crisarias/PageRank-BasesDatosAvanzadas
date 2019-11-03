import os

import random
file = open("./filename.txt", "w")
CantNodos = random.randrange(100)
print(CantNodos)
file.write(str(CantNodos)+os.linesep)
for lineas in range(CantNodos):
    vertex =  random.randrange(CantNodos)
    
    print(lineas,end = '')
    print(" ",end = '')
    file.write(str(lineas) + " ")
    
    for x in range(vertex):
        linkto= random.randrange(CantNodos)
        file.write(str(linkto)+ " ")
        print(linkto,end = '')
        print(" ",end = '')
    print("")
    file.write(os.linesep)

file.close()