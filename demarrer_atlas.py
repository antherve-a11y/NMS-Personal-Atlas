"""Atlas Windows : un double-clic, installation automatique une seule fois.

Runtime Python/SQLite embarque. Dependances telechargees dans runtime/site-packages
au premier demarrage seulement. Aucune modification de l'installation Windows.
Les donnees personnelles restent toujours dans <racine Atlas>/NMSFaunaAnalyzerData.
"""
from __future__ import annotations

import hashlib
import importlib.util
import json
import os
from pathlib import Path
import shutil
import socket
import subprocess
import sys
import tempfile
import threading
import time
import urllib.request
import webbrowser

ROOT = Path(__file__).resolve().parents[2]
DATA = ROOT / 'NMSFaunaAnalyzerData'
RUNTIME = ROOT / 'runtime'
PYTHON = RUNTIME / 'python' / 'python.exe'
SITE = RUNTIME / 'site-packages'
READY = RUNTIME / 'dependencies-ready.json'
REQ_PYPI = ROOT / 'runtime' / 'launcher' / 'requirements_atlas_pypi.txt'
REQ_NAMEGEN = ROOT / 'runtime' / 'launcher' / 'requirements_atlas_namegen.txt'
VENDOR_NAMEGEN = ROOT / 'runtime' / 'launcher' / 'vendor' / 'nms_namegen'
SCHEMA = 'atlas-private-runtime-v2-bundled-namegen'
PIP_BOOTSTRAP = 'https://bootstrap.pypa.io/get-pip.py'
MODULES = ('streamlit', 'pandas', 'numpy', 'onnxruntime', 'rapidocr',
           'lz4', 'PIL', 'xxhash', 'mss', 'nms_namegen.planet',
           'nms_namegen.system')


def _log(message: str) -> None:
    print('[ATLAS] ' + message, flush=True)


def _ui(step: str, message: str) -> None:
    """Evens GUI sans changer les logs ni les integrations existantes."""
    print(f'@ATLAS_UI\t{step}\t{message}', flush=True)


def _fingerprint() -> str:
    h = hashlib.sha256()
    for path in (REQ_PYPI, REQ_NAMEGEN):
        h.update(path.read_bytes())
        h.update(b'\0')
    # Versionner egalement les fichiers fournis du generateur tiers.
    # Une modification de ce module declenchera une installation de maintenance.
    for path in sorted(VENDOR_NAMEGEN.rglob('*')):
        if path.is_file():
            h.update(str(path.relative_to(VENDOR_NAMEGEN)).encode('utf-8'))
            h.update(path.read_bytes())
    h.update(SCHEMA.encode('ascii'))
    return h.hexdigest()


def _installed_ok() -> bool:
    try:
        marker = json.loads(READY.read_text(encoding='utf-8'))
        if marker.get('schema') != SCHEMA or marker.get('requirements_sha256') != _fingerprint():
            return False
        if marker.get('python_version') != list(sys.version_info[:3]):
            return False
        if not SITE.is_dir():
            return False
        # Peu couteux : ne charge pas NumPy, Streamlit ou ONNX a chaque lancement.
        return all(importlib.util.find_spec(module) is not None for module in MODULES)
    except (OSError, ValueError, ImportError, ModuleNotFoundError, AttributeError):
        return False


def _download_get_pip() -> Path:
    dest = RUNTIME / 'bootstrap' / 'get-pip.py'
    if dest.is_file() and dest.stat().st_size > 500_000:
        return dest
    dest.parent.mkdir(parents=True, exist_ok=True)
    _log('Telechargement securise du gestionnaire de bibliotheques depuis PyPA...')
    temp = dest.with_suffix('.download')
    try:
        req = urllib.request.Request(PIP_BOOTSTRAP, headers={'User-Agent':'NMS-Personal-Atlas/12.7'})
        with urllib.request.urlopen(req, timeout=45) as response, temp.open('wb') as fd:
            if response.geturl().split('/')[2] != 'bootstrap.pypa.io':
                raise RuntimeError('Redirection inattendue pour le gestionnaire pip.')
            size = 0
            while True:
                chunk = response.read(1 << 19)
                if not chunk:
                    break
                size += len(chunk)
                if size > 10 * 1024 * 1024:
                    raise RuntimeError('Fichier get-pip.py anormalement grand.')
                fd.write(chunk)
        if not (500_000 < size < 10 * 1024 * 1024):
            raise RuntimeError('Telechargement incomplet de pip.')
        os.replace(temp, dest)
    finally:
        temp.unlink(missing_ok=True)
    return dest


def _run(args: list[str], log_file: Path, env: dict[str, str]) -> None:
    with log_file.open('a', encoding='utf-8') as log:
        log.write('\n=== Atlas dependency step ===\n')
        log.flush()
        # Montrer l'avancement pendant l'installation, et conserver un journal.
        process = subprocess.Popen(args, cwd=ROOT, env=env, stdout=subprocess.PIPE,
                                   stderr=subprocess.STDOUT, text=True,
                                   encoding='utf-8', errors='replace', bufsize=1,
                                   creationflags=getattr(subprocess, 'CREATE_NO_WINDOW', 0))
        assert process.stdout is not None
        for line in process.stdout:
            log.write(line)
            log.flush()
            print('  ' + line.rstrip(), flush=True)
        result = process.wait()
        if result:
            raise RuntimeError(f'Installation echouee (code {result}). Voir {log_file}')


def _ensure_pip(log_file: Path, env: dict[str, str]) -> None:
    if importlib.util.find_spec('pip') is not None:
        return
    pip_script = _download_get_pip()
    _log('Preparation de pip dans le Python prive Atlas...')
    _run([str(PYTHON), str(pip_script), '--no-warn-script-location',
          '--disable-pip-version-check', '--no-input'], log_file, env)
    # Le subprocess a ajoute pip dans Lib/site-packages ; le prochain subprocess le voit.


def _include_bundled_namegen(stage: Path) -> None:
    """Copier le module MIT deja inclus : aucun build Hatchling ni acces GitHub."""
    required = ('__init__.py', 'generator.py', 'prng.py', 'iprng.py',
                'planet.py', 'region.py', 'system.py', 'roman.py', 'letter_map.json')
    missing = [name for name in required if not (VENDOR_NAMEGEN / name).is_file()]
    if missing:
        raise RuntimeError('Le module NMS NameGen integre est incomplet : ' + ', '.join(missing))
    dest = stage / 'nms_namegen'
    if dest.exists():
        shutil.rmtree(dest)
    shutil.copytree(VENDOR_NAMEGEN, dest)
    _log('NMS NameGen 2.0 integre depuis Atlas (aucun telechargement GitHub ni Hatchling).')


def _verify_stage(stage: Path, log_file: Path, env: dict[str, str]) -> None:
    test = (
        'import sys,pathlib; target=pathlib.Path(sys.argv[1]).resolve(); '
        'sys.path.insert(0, str(target)); '
        'import streamlit,pandas,numpy,sqlite3,onnxruntime,rapidocr,lz4,PIL,xxhash,mss; '
        'import nms_namegen.planet,nms_namegen.system; '
        'import atlas_core.database; '
        'mods=(streamlit,pandas,numpy,onnxruntime,rapidocr,lz4,PIL,xxhash,mss,nms_namegen.planet); '
        'assert all(pathlib.Path(m.__file__).resolve().is_relative_to(target) for m in mods), '
        '"Modules non charges depuis la nouvelle installation"; '
        'print("ATLAS: import smoke test OK")'
    )
    _run([str(PYTHON), '-B', '-c', test, str(stage)], log_file, env)


def _install_once(log_file: Path) -> None:
    if _installed_ok():
        _ui('VERIFY', 'Bibliothèques déjà installées : démarrage direct')
        return
    if sys.platform != 'win32' or sys.version_info[:2] != (3, 13):
        raise RuntimeError('Atlas requiert son Python Windows 3.13 embarque.')

    _ui('INSTALL', 'Téléchargement et installation des bibliothèques (première utilisation)')
    _log('Premiere utilisation : installation automatique des bibliotheques. Aucun clic supplementaire necessaire.')
    RUNTIME.mkdir(parents=True, exist_ok=True)
    # Les paquets et leurs caches restent a l'interieur du dossier Atlas.
    # Ne pas laisser un pip.cfg personnel ou PIP_* systeme modifier l'installation.
    env = {key: value for key, value in os.environ.items() if not key.startswith('PIP_')}
    env['PYTHONNOUSERSITE'] = '1'
    env['PIP_NO_INPUT'] = '1'
    env['PIP_DISABLE_PIP_VERSION_CHECK'] = '1'
    env['PIP_CACHE_DIR'] = str(DATA / 'runtime' / 'pip-cache')
    # IMPORTANT Windows : pip place normalement ses fichiers temporaires sur C:
    # meme si Atlas est sur D:. Un deplacement inter-disques provoquait
    # WinError 17 puis WinError 3 pour les chemins profonds d'ONNX Runtime.
    # Le repertoire de travail de pip reste donc sur le MEME disque qu'Atlas.
    pip_temp = RUNTIME / '.tmp'
    pip_temp.mkdir(parents=True, exist_ok=True)
    env['TMP'] = str(pip_temp)
    env['TEMP'] = str(pip_temp)
    env['TMPDIR'] = str(pip_temp)
    # Ne pas creer les chemins __pycache__ profonds au cours de l'installation.
    env['PYTHONDONTWRITEBYTECODE'] = '1'
    env['PIP_INDEX_URL'] = 'https://pypi.org/simple'
    env['PIP_CONFIG_FILE'] = os.devnull
    env.pop('PYTHONHOME', None)
    env.pop('PYTHONPATH', None)
    Path(env['PIP_CACHE_DIR']).mkdir(parents=True, exist_ok=True)
    _ensure_pip(log_file, env)
    # Prefixe court pour rester sous la limite de longueur des chemins Windows.
    stage_dir = Path(tempfile.mkdtemp(prefix='.a-', dir=RUNTIME))
    staged_packages = stage_dir / 'site-packages'
    staged_packages.mkdir()
    backup = RUNTIME / '.previous-site-packages'
    try:
        _log('Installation de Streamlit, pandas, NumPy, reconnaissance et compression...')
        _run([str(PYTHON), '-m', 'pip', 'install', '--disable-pip-version-check',
              '--no-input', '--only-binary=:all:', '--ignore-installed',
              '--no-compile', '--progress-bar', 'off', '--target', str(staged_packages), '-r', str(REQ_PYPI)], log_file, env)
        # Module tiers fourni dans le ZIP de distribution : pas de compilation
        # Hatchling et aucun telechargement direct depuis GitHub.
        _include_bundled_namegen(staged_packages)
        _ui('VERIFY', 'Vérification des modules et bibliothèques Windows')
        _log('Verification des DLL Windows et des modules Atlas...')
        _verify_stage(staged_packages, log_file, env)
        # Publication atomique autant que possible ; en cas d'erreur, restaurer.
        if backup.exists():
            shutil.rmtree(backup)
        had_previous = SITE.is_dir()
        if had_previous:
            SITE.rename(backup)
        try:
            staged_packages.rename(SITE)
        except Exception:
            if had_previous and backup.is_dir() and not SITE.exists():
                backup.rename(SITE)
            raise
        marker = {
            'schema': SCHEMA,
            'requirements_sha256': _fingerprint(),
            'python_version': list(sys.version_info[:3]),
            'status': 'imports-verified',
        }
        temp_marker = READY.with_suffix('.tmp')
        temp_marker.write_text(json.dumps(marker, indent=2), encoding='utf-8')
        os.replace(temp_marker, READY)
        _log('Bibliotheques installees et verifiees. Les prochains lancements les reutiliseront.')
        if backup.exists():
            shutil.rmtree(backup, ignore_errors=True)
    finally:
        shutil.rmtree(stage_dir, ignore_errors=True)


def _port_available(port: int) -> bool:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        return sock.connect_ex(('127.0.0.1', port)) != 0


def main() -> int:
    if sys.platform != 'win32' or sys.version_info[:2] != (3, 13):
        _log('Python 3.13 Windows 64 bits necessaire.')
        return 3
    if Path(sys.executable).resolve() != PYTHON.resolve():
        _log('Le programme doit utiliser le Python embarque dans Atlas.')
        return 4
    if not (RUNTIME / 'python' / 'python313._pth').is_file():
        _log('Fichiers du Python embarque incomplets. Reextraire Atlas.')
        return 5

    os.environ['NMS_FAUNA_DATA_DIR'] = str(DATA)
    os.environ['NMS_RUNTIME_DIR'] = str(DATA / 'runtime')
    os.environ['STREAMLIT_BROWSER_GATHER_USAGE_STATS'] = 'false'
    os.environ['STREAMLIT_SERVER_HEADLESS'] = 'true'
    os.environ['PYTHONNOUSERSITE'] = '1'
    os.environ['PYTHONDONTWRITEBYTECODE'] = '1'
    os.environ['PYTHONHOME'] = ''
    os.environ['PYTHONPATH'] = ''
    os.environ['ATLAS_OFFLINE_ONLY'] = '1'  # Pas de telechargements de donnees annexes inattendus.
    sys.path.insert(0, str(ROOT))

    # Do not perform a database-directory swap while a previous Atlas instance
    # is still serving requests on the configured local port.
    if (ROOT / '.atlas-manual-recovery-pending.json').is_file():
        _recovery_port = int(os.environ.get('ATLAS_RUN_PORT', '8532'))
        if not 1024 <= _recovery_port <= 65535 or not _port_available(_recovery_port):
            _ui('ERROR', 'Fermez toutes les fenetres Atlas avant la recuperation des donnees.')
            _log('Recuperation manuelle refusee : session Atlas encore ouverte.')
            return 8

    _ui('DATA', 'Préparation des données personnelles')
    from atlas_core.storage_bootstrap_offline import bootstrap_local_data, apply_pending_manual_recovery
    try:
        _manual_recovery = apply_pending_manual_recovery(ROOT)
        if _manual_recovery:
            _log('Recuperation manuelle : ' + _manual_recovery)
        outcome = bootstrap_local_data(ROOT)
    except Exception as exc:
        _ui('ERROR', f'Erreur d’initialisation du stockage : {exc}')
        _log(f'Initialisation des donnees interrompue pour les proteger : {exc}')
        return 6
    _log(f'Donnees locales : {DATA} ({outcome})')

    log_dir = DATA / 'runtime' / 'logs'
    log_dir.mkdir(parents=True, exist_ok=True)
    install_log = log_dir / 'installation_dependances.log'
    _ui('CHECK', 'Vérification des composants existants')
    try:
        _install_once(install_log)
    except Exception as exc:
        _ui('ERROR', f'Installation impossible : {exc}')
        _log(f'Installation automatique incomplete : {exc}')
        _log('Verifier la connexion Internet ; relancer le meme LANCER_ATLAS.bat pour reprendre.')
        return 7

    port = int(os.environ.get('ATLAS_RUN_PORT', '8532'))
    if not 1024 <= port <= 65535 or not _port_available(port):
        _log(f'Port local {port} indisponible (Atlas est peut-etre deja lance).')
        return 8

    local_runtime = DATA / 'runtime'
    (local_runtime / 'temp').mkdir(parents=True, exist_ok=True)
    os.environ['TMP'] = str(local_runtime / 'temp')
    os.environ['TEMP'] = str(local_runtime / 'temp')
    os.environ.pop('PYTHONPYCACHEPREFIX', None)  # aucune generation de .pyc longs sur Windows
    command = [str(PYTHON), '-B', '-m', 'streamlit', 'run', str(ROOT / 'app.py'),
               '--server.address', '127.0.0.1', '--server.port', str(port),
               '--server.headless', 'true', '--server.fileWatcherType', 'none',
               '--browser.gatherUsageStats', 'false']

    def open_browser() -> None:
        for _ in range(90):
            time.sleep(0.35)
            if not _port_available(port):
                webbrowser.open(f'http://127.0.0.1:{port}')
                return

    threading.Thread(target=open_browser, daemon=True).start()
    _log(f'Demarrage local : http://127.0.0.1:{port}')
    _ui('START', 'Ouverture de votre Atlas')
    try:
        return subprocess.call(command, cwd=ROOT, env=os.environ.copy(),
                               creationflags=getattr(subprocess, 'CREATE_NO_WINDOW', 0))
    except KeyboardInterrupt:
        return 0

if __name__ == '__main__':
    raise SystemExit(main())
