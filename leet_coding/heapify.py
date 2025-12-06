def heapify(nums, n, i):
    largest = i
    left = 2 * i + 1
    right = 2 * i + 2
    if left < n and nums[largest] > nums[left]:
        # violation 
        largest = left
    if right < n and nums[largest] > nums[right]:
        largest = right
    
    if largest != i:
        nums[largest], nums[i] = nums[i], nums[largest]
        heapify(nums, n, largest)


def insert_element(nums, val):
    nums.append(val)
    index_of_insertion = len(nums)-1
    i = index_of_insertion
    while i >0 :
        parent = (i-1)//2
        if nums[parent] >  nums[i]:
            nums[parent], nums[i] = nums[i], nums[parent]
            i = parent
        else:
            break

def get_max():
    if len(nums) == 0:
        return None
    root = nums[0]
    nums[0] = nums[-1]
    nums.pop()
    heapify(nums, len(nums), 0)
    return root


nums = [3, 9, 2, 1, 4, 5, 4, 2]



print(nums)
# insert_element(nums, 30)
# print(nums)
insert_element(nums, 0)
print(nums)