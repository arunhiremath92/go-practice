def zig_zag(s, max_rows):
    if max_rows<2:
        return s
    i = 0
    index = 0
    a = []
    str_len = len(s)
    while index < len(s):
        a.append([])
        for r in range(0, max_rows):
            a[i].append(-1)

        if i % (max_rows - 1) == 0:
            j = 0
            while j < (max_rows) and index < str_len:
                a[i][j] = s[index]
                index += 1
                j = j + 1
            j -= 1
        
        if i % (max_rows - 1) != 0 and index < str_len:
            j = j - 1
            a[i][j] = s[index]
            index += 1
        # prepare for next iteration
        i += 1
    # construct the final string
    result = ""
    for i in range(0, max_rows):
        for j in range(0, len(a)):
            if a[j][i] != -1:
                result = result + a[j][i]
    return result
# 
print(zig_zag("a",1))
print(zig_zag("ab",1))
print(zig_zag("abc",2))
print(zig_zag("paypalishiring",2))
print(zig_zag("paypalishiring",3))
print(zig_zag("paypalishiring",4))



# def zig_zag_vertical(s, max_rows):
#     i = 0
#     index = 0
#     a = []
#     str_len = len(s)
#     while index < len(s):
#         a.append([])
#         for r in range(0, max_rows):
#             a[i].append(-1)
#         if i % (max_rows - 1) == 0:
#             j = max_rows - 1
#             while j >=0  and index < str_len:
#                 a[i][j] = s[index]
#                 index += 1
#                 j = j - 1
#             j += 1
        
#         if i % (max_rows - 1) != 0 and index < str_len:
#             j = j + 1
#             a[i][j] = s[index]
#             index += 1
#         i += 1
