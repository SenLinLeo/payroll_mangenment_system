====
什么东西？What?
----
薪资核算系统使用 golang开发，基于beego，页面基于layUi,目前版本1.0.0   
此系统因涉及到员工个人信息、薪资信息，需保障数据信息安全，做好敏感信息防护，数据加密、登录权限控制等开发自主完成；
统计汇总、导入、计算等功能都涉及到系统的性能，在保证系统数据的实时性同时也要保证系统的稳定性。

有什么价值？
----
1、RBAC权限完善，多角色管理系统    
2、后台界面完整，多标签页面    
3、API相关页面有比较复杂的使用案例    
所以，可以作为一个基础框架使用，快速开发。初学者还可以作为熟悉beego使用。 
4、可用于中小企业薪资核算系统使用  


用到了哪些？
----
1、界面框架layUI2.4.5    
2、makedown.md   
3、beego1.8
4、Ztree   

安装方法    
----
1、go get xxx    
2、创建mysql数据库，并将geek-nebula-api-admin.sql导入    
3、修改config 配置数据库    
4、运行 go build    
5、运行 ./run.sh start|stop
6. 访问 http://127.0.0.1:8082/home
 
用户名：admin 密码：george518 


