from setuptools import setup, find_packages

setup(
    name='l10n',
    version='0.0.1',
    packages=find_packages(),
    install_requires=[
        'pyyaml',
    ],
    description='Localization Generator Tool',
    license='Apache License 2.0',
    url='https://github.com/A1ex3/l10n',
)