import os

def rename_files_in_directory(directory, file_extension, string_to_remove):
    # ディレクトリ内のすべてのファイルをチェック
    for filename in os.listdir(directory):
        # 指定した拡張子を持つファイルか確認
        if filename.endswith(file_extension):
            # ファイル名から指定された文字列を削除
            new_filename = filename.replace(string_to_remove, "")
            # 元のファイルパスと新しいファイルパスを取得
            old_path = os.path.join(directory, filename)
            new_path = os.path.join(directory, new_filename)
            # ファイル名を変更
            os.rename(old_path, new_path)
            print(f"Renamed: '{filename}' -> '{new_filename}'")

# パラメータの設定
target_directory = "tachibana/business_functions/response"  # 対象のディレクトリを指定
file_ext = ".go"                           # 対象とするファイルの拡張子
string_to_remove = "res_"            # ファイル名から削除する文字列

# 関数の呼び出し
rename_files_in_directory(target_directory, file_ext, string_to_remove)
