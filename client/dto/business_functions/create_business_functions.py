import os

# ベースディレクトリ名
base_dir = "business_functions"

# ファイルリスト
files = [
    "req_new_order.go", "res_new_order.go",
    "req_correct_order.go", "res_correct_order.go",
    "req_cancel_order.go", "res_cancel_order.go",
    "req_cancel_order_all.go", "res_cancel_order_all.go",
    "req_genbutu_kabu_list.go", "res_genbutu_kabu_list.go",
    "req_shinyou_tategyoku_list.go", "res_shinyou_tategyoku_list.go",
    "req_zan_kai_kanougaku.go", "res_zan_kai_kanougaku.go",
    "req_zan_shinki_kano_ijiritu.go", "res_zan_shinki_kano_ijiritu.go",
    "req_zan_uri_kanousuu.go", "res_zan_uri_kanousuu.go",
    "req_order_list.go", "res_order_list.go",
    "req_order_list_detail.go", "res_order_list_detail.go",
    "req_zan_kai_summary.go", "res_zan_kai_summary.go",
    "req_zan_kai_kanougaku_suii.go", "res_zan_kai_kanougaku_suii.go",
    "req_zan_kai_genbutu_kaituke_syousai.go", "res_zan_kai_genbutu_kaituke_syousai.go",
    "req_zan_kai_sinyou_sinkidate_syousai.go", "res_zan_kai_sinyou_sinkidate_syousai.go",
    "req_zan_real_hosyoukin_ritu.go", "res_zan_real_hosyoukin_ritu.go"
]

# ディレクトリの作成
os.makedirs(base_dir, exist_ok=True)

# ファイルの作成
for file in files:
    file_path = os.path.join(base_dir, file)

print(f"'{base_dir}' ディレクトリとファイル群が作成されました！")
