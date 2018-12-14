input_ = """0 pwm 12
wait 50
0 pwm 25
wait 50
0 pwm 50
1 pwm 12
wait 50
0 pwm 100
1 pwm 25
wait 50
0 pwm 200
1 pwm 50
2 pwm 12
wait 50
1 pwm 100
2 pwm 25
wait 50
0 pwm 0
1 pwm 200
2 pwm 50
3 pwm 12
wait 50
2 pwm 100
3 pwm 25
wait 50
1 pwm 0
2 pwm 200
3 pwm 50
wait 50
3 pwm 100
wait 50
3 pwm 200
2 pwm 0
wait 100
3 pwm 0""".split("\n")
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
        
    
    
