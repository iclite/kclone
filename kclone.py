#!/usr/bin/env python3

from __future__ import annotations
import argparse
import os
import json
import shutil
import time
import subprocess
import git
from git.exc import GitCommandError
from git import Repo
from urllib.parse import urlparse

from rich import console, progress

# 使用 USERPROFILE 环境变量来获取用户主目录
USER_HOME = os.environ.get('USERPROFILE') or os.path.expanduser('~')
CONFIG_FILE = os.path.join(USER_HOME, '.kclone_config.json')
DEFAULT_CLONE_DIR = os.path.join(USER_HOME, 'gitworks')




class GitRemoteProgress(git.RemoteProgress):
    OP_CODES = [
        "BEGIN",
        "CHECKING_OUT",
        "COMPRESSING",
        "COUNTING",
        "END",
        "FINDING_SOURCES",
        "RECEIVING",
        "RESOLVING",
        "WRITING",
    ]
    OP_CODE_MAP = {
        getattr(git.RemoteProgress, _op_code): _op_code for _op_code in OP_CODES
    }

    def __init__(self) -> None:
        super().__init__()
        self.progressbar = progress.Progress(
            progress.SpinnerColumn(),
            # *progress.Progress.get_default_columns(),
            progress.TextColumn("[progress.description]{task.description}"),
            progress.BarColumn(),
            progress.TextColumn("[progress.percentage]{task.percentage:>3.0f}%"),
            "eta",
            progress.TimeRemainingColumn(),
            progress.TextColumn("{task.fields[message]}"),
            console=console.Console(),
            transient=False,
        )
        self.progressbar.start()
        self.active_task = None

    def __del__(self) -> None:
        # logger.info("Destroying bar...")
        self.progressbar.stop()

    @classmethod
    def get_curr_op(cls, op_code: int) -> str:
        """Get OP name from OP code."""
        # Remove BEGIN- and END-flag and get op name
        op_code_masked = op_code & cls.OP_MASK
        return cls.OP_CODE_MAP.get(op_code_masked, "?").title()

    def update(
        self,
        op_code: int,
        cur_count: str | float,
        max_count: str | float | None = None,
        message: str | None = "",
    ) -> None:
        # Start new bar on each BEGIN-flag
        if op_code & self.BEGIN:
            self.curr_op = self.get_curr_op(op_code)
            # logger.info("Next: %s", self.curr_op)
            self.active_task = self.progressbar.add_task(
                description=self.curr_op,
                total=max_count,
                message=message,
            )

        self.progressbar.update(
            task_id=self.active_task,
            completed=cur_count,
            message=message,
        )

        # End progress monitoring on each END-flag
        if op_code & self.END:
            # logger.info("Done: %s", self.curr_op)
            self.progressbar.update(
                task_id=self.active_task,
                message=f"[bright_black]{message}",
            )

def load_config():
    """加载配置文件"""
    if os.path.exists(CONFIG_FILE):
        with open(CONFIG_FILE, 'r') as f:
            return json.load(f)
    return {"default_clone_dir": DEFAULT_CLONE_DIR}

def save_config(config):
    """保存配置文件"""
    with open(CONFIG_FILE, 'w') as f:
        json.dump(config, f, indent=2)

def parse_git_url(url):
    """解析Git URL,返回主机名、用户名和仓库名"""
    parsed_url = urlparse(url)
    path_parts = parsed_url.path.strip('/').split('/')
    return parsed_url.hostname, path_parts[0], path_parts[1].rstrip('.git')


def safe_remove_directory(path):
    """安全地删除目录，处理权限错误"""
    max_attempts = 3
    for attempt in range(max_attempts):
        try:
            shutil.rmtree(path)
            return True
        except PermissionError:
            if attempt < max_attempts - 1:
                print(f"删除目录失败，正在重试... (尝试 {attempt + 1}/{max_attempts})")
                time.sleep(1)  # 等待1秒后重试
                # 尝试使用系统命令强制删除
                try:
                    if os.name == 'nt':  # Windows
                        subprocess.run(['rd', '/s', '/q', path], check=True, shell=True)
                    else:  # Unix-like
                        subprocess.run(['rm', '-rf', path], check=True)
                    return True
                except subprocess.CalledProcessError:
                    continue
            else:
                print(f"无法删除目录 {path}。请手动删除后重试。")
                return False

def kclone(url, base_dir):
    """克隆仓库到指定的目录结构中"""
    hostname, username, repo_name = parse_git_url(url)
    
    # 创建目录结构
    full_path = os.path.join(base_dir, hostname, username)
    os.makedirs(full_path, exist_ok=True)
    
    # 克隆仓库
    clone_path = os.path.join(full_path, repo_name)
    if os.path.exists(clone_path) and os.listdir(clone_path):
        print(f"目标目录 '{clone_path}' 已存在且不为空。")
        choice = input("请选择操作：[O]覆盖 / [R]重命名 / [C]取消 ").lower()
        
        if choice == 'o':
            if not safe_remove_directory(clone_path):
                return
        elif choice == 'r':
            new_name = input("请输入新的目录名：")
            clone_path = os.path.join(os.path.dirname(clone_path), new_name)
        else:
            print("操作已取消。")
            return

    try:
        Repo.clone_from(url, clone_path, recurse_submodules=True, progress=GitRemoteProgress())
        print(f"仓库已成功克隆到 {clone_path}")
        print("")
        print(f"    explorer {clone_path}")
        print(f"    code     {clone_path}")
        print("")
                
    except GitCommandError as e:
        print(f"克隆仓库时出错：{e}")

def main():
    config = load_config()
    
    parser = argparse.ArgumentParser(description="克隆Git仓库到组织化的目录结构中")
    parser.add_argument("url", nargs="?", help="要克隆的Git仓库URL")
    parser.add_argument("-d", "--directory", help="基础克隆目录")
    parser.add_argument("--set-default", help="设置新的默认克隆目录")
    
    args = parser.parse_args()

    if args.set_default:
        config["default_clone_dir"] = os.path.abspath(os.path.expanduser(args.set_default))
        save_config(config)
        print(f"默认克隆目录已设置为: {config['default_clone_dir']}")
        return

    if not args.url:
        print(f"当前默认克隆目录: {config['default_clone_dir']}")
        return

    base_dir = args.directory or config["default_clone_dir"]
    kclone(args.url, base_dir)

if __name__ == "__main__":
    main()