# print(not (not True))
print(True)
# print(not (not True and False))
print(not (not True and False))
# print(not not True)
print(True)
# print(not not False)
print(False)
# print(not not not True)
print(not True)
# print(not not not not True)
print(True)
# print(not (((not not not True))))
print(True)
# print(not (not not True and False))
print(not (True and False))


# if (not (not testWhatever(12))):
#     return "Hello"
if (testWhatever(12)):
    return "Hello"
