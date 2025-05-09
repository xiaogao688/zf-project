#!/bin/bash
echo -e "\033[31m 脚本使用场景：仅限从Vcenter虚拟化环境模板部署的虚拟机根目录扩容，扩容逻辑卷 centos-root，参数已写死，其他场景不用 \033[0m"
echo -e "
\033[31m 查看硬盘：\033[0m
$(lsblk)
------------------------
\033[31m 查看当前的VG: \033[0m
$(vgs)
------------------------
\033[31m 查看前存在LV分区如下：\033[0m
$(lvs)
"
##################   用户输入  ################

input(){
echo -e "
如果输入错误字符或者闪跳，请Ctrl +c  退出重新输入
根据上述输出，请输入你要扩容到根目录的磁盘; 磁盘格式为：\033[31m sdb 或 sda 或sdc等格式 \033[0m "
read -p "请输入要分区的磁盘:"  disk  ;
}
echo -e "\033[31m 当前是新建LVM卷和挂载新分区 \033[0m"

input ;
echo -e "\033[31m 即将$disk磁盘加入到根目录centos-root逻辑卷中\033[0m"

echo -e "\033[31m 创建pv \033[0m"
pvcreate /dev/$disk

echo -e "\033[31m 查看当前pv \033[0m"
pvs

echo -e "\033[31m 扩容硬盘到vg卷组\033[0m"
vgextend centos /dev/$disk

echo -e "\033[31m 查看当前vg卷组 \033[0m"
vgs

echo -e "\033[31m 扩容当前所有空间至centos-root逻辑卷中 \033[0m"
lvextend -l +100%FREE /dev/mapper/centos-root
xfs_growfs /dev/mapper/centos-root

echo -e "\033[31m 查看当前挂载目录，关注根目录是否已经扩容\033[0m"
df -h
