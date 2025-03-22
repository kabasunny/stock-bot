import os

# ディレクトリとファイルの構造を定義
structure = {
        "cmd": {
            "myapp": ["main.go"]
        },
        "internal": {
            "app": {
                "swingtrade": ["usecase.go", "service.go"],  # service.go は必要に応じて
                "daytrade": ["usecase.go", "service.go"]  # service.go は必要に応じて
            },
            "infrastructure": {
                "client": ["client.go", "mock_client.go"],
                "repository": ["repository.go", "mock_repository.go"],
                "eventhandler": ["eventhandler.go"]
            },
            "interface": {  # プレゼンテーション層
                "web": {  # Web UI用のパッケージ
                    "templates": [],  # HTMLテンプレート用
                    "handlers.go": None
                },
                "agents": {
                    "swingtrade": ["swingtrade.go"],
                    "daytrade": ["daytrade.go"]
                }
            },
        },
        "pkg": {},
        "domain": {
            "model": ["order.go", "position.go"],  # 他のファイルも追加可能
            "repository": ["order_repository.go"],
            "service": [],
            "event": []
        
    }
}

# ディレクトリとファイルを作成する関数
def create_structure(base_path, structure):
    for name, content in structure.items():
        path = os.path.join(base_path, name)
        if isinstance(content, dict):
            os.makedirs(path, exist_ok=True)
            create_structure(path, content)
        elif isinstance(content, list):
            os.makedirs(path, exist_ok=True)
            for file_name in content:
                open(os.path.join(path, file_name), 'w').close()
        elif content is None:
            open(path, 'w').close()

# 実行
base_path = "./"  # ベースとなるディレクトリ
create_structure(base_path, structure)
