import os
import sys

# adjust path so underlying library files can access each other without any hassle
# current path is python3/ so we need to descend once
dir_path = os.path.dirname(os.path.realpath(__file__))
sys.path.append(os.path.abspath(os.path.join(dir_path, "..")))
