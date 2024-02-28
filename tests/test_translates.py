import json
import pytest
from tests.app_localization import AppLocalization

@pytest.fixture
def localization(request):
    locale = request.param
    return AppLocalization(locale)

@pytest.fixture
def json_data(request):
    locale = request.param
    with open(f'tests/translates/l10n_{locale}.json', encoding='utf-8') as f:
        return json.load(f)

locales = ['en', 'ru']

@pytest.mark.parametrize('localization,json_data', [(loc, loc) for loc in locales], indirect=True)
def test_hello_world(localization, json_data):
    assert localization.of().helloWorld == json_data['helloWorld']

@pytest.mark.parametrize('localization,json_data', [(loc, loc) for loc in locales], indirect=True)
def test_bye(localization, json_data):
    expected_bye = json_data['bye'].replace('{value}', json_data['#bye']['variables']['value']['defaultValue'])
    assert localization.of().bye() == expected_bye

@pytest.mark.parametrize('localization,json_data', [(loc, loc) for loc in locales], indirect=True)
def test_number_of_users(localization, json_data):
    locale = localization.current_locale
    if 'numberOfUsers' in json_data:
        expected_number_of_users = json_data['numberOfUsers'].replace('{number}', str(json_data['#numberOfUsers']['variables']['number']['defaultValue']))
        assert localization.of().numberOfUsers(json_data['#numberOfUsers']['variables']['values']['type']) == expected_number_of_users
