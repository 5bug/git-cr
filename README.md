# git-cr
实现gitlab下主干开发模式的CR流程

## git-cr流程
这里以master分支为例，来说下具体的流程：  
1.master分支设置为保护模式：不允许任何人提交，需要通过MR的方式合入代码.   
2.大家都在本地都在master分支上写代码，不需要拉出feature分支.   
3.写完代码后，直接git push推送是不成功的，这个时候需要使用git cr命令来替代git push.   
4.git cr命令执行后会先同步远端master代码带本地，若有冲突在本地解决，无冲突则将本地代码推送到远端临时分支，同时发起一个MR合并请求.   
5.在MR里做CR，CR通过后即可合入master.  

## 安装方法
go install github.com/5bug/git-cr@latest

## 使用方法
git add xxx   
git commit -m "feat: xxx"  
git cr
