input_ = """0 pwm 0
1 pwm 0
2 pwm 0
3 pwm 0

0 pwm 255
wait 255
wait 255
wait 255
wait 255

1 pwm 255
0 pwm 0
wait 255
wait 255
wait 255
wait 255

2 pwm 255
1 pwm 0
wait 255
wait 255
wait 255
wait 255


3 pwm 255
2 pwm 0
wait 255
wait 255
wait 255
wait 255

2 pwm 255
3 pwm 0
wait 255
wait 255
wait 255
wait 255

1 pwm 255
2 pwm 0
wait 255
wait 255
wait 255
wait 255

0 pwm 255
1 pwm 0
wait 255
wait 255
wait 255
wait 255
""".split("\n")
def helper(a) :
    if a=="pwm":
        return "true"
    else:
        return "false"
for line in input_:
    tokens = line.split(" ")
    ##print(line)
    if len(tokens) > 1:
        if len(tokens) == 3:
            pin = tokens[0]
            pwm = helper(tokens[1])
            byte = tokens[2]
            print("command{pin: led["+pin+"], pwm: "+pwm+", instruction: byte("+byte+")},",end='')
        elif len(tokens) == 2:
            pin = "nil"
            pwm = "true"
            byte = tokens[1]
            print("command{pin: nil, pwm: "+pwm+", instruction: byte("+byte+")},",end='')
print()
        
    
    
