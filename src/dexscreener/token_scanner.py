import asyncio
import aiohttp
from datetime import datetime
import pandas as pd
import json
from typing import Dict, List
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class DexTokenScanner:
    def __init__(self):
        self.base_url = "https://api.dexscreener.com"
        self.headers = {
            "User-Agent": "Mozilla/5.0",
            "Accept": "application/json"
        }
        self.session = None
        self.analyzed_tokens = set()
        self.promising_tokens = []
        
        # Rate limiting - 60 requests per minute for profiles
        self.profile_rate_limit = 60
        self.profile_interval = 60 / self.profile_rate_limit

    async def initialize(self):
        """Initialize aiohttp session"""
        self.session = aiohttp.ClientSession(headers=self.headers)

    async def close(self):
        """Close aiohttp session"""
        if self.session:
            await self.session.close()

    async def get_latest_token_profiles(self) -> List[Dict]:
        """Fetch latest token profiles"""
        try:
            endpoint = f"{self.base_url}/token-profiles/latest/v1"
            
            async with self.session.get(endpoint) as response:
                if response.status != 200:
                    logger.error(f"Error fetching token profiles: Status {response.status}")
                    return []
                
                profiles = await response.json()
                logger.info(f"Retrieved {len(profiles)} token profiles")
                return profiles

        except Exception as e:
            logger.error(f"Error fetching token profiles: {str(e)}")
            return []

    async def get_pair_info(self, chain_id: str, pair_address: str) -> Dict:
        """Get detailed information about a specific pair"""
        try:
            endpoint = f"{self.base_url}/latest/dex/pairs/{chain_id}/{pair_address}"
            
            async with self.session.get(endpoint) as response:
                if response.status != 200:
                    logger.error(f"Error fetching pair info: Status {response.status}")
                    return {}
                
                data = await response.json()
                return data.get('pairs', [{}])[0]

        except Exception as e:
            logger.error(f"Error fetching pair info: {str(e)}")
            return {}

    async def analyze_token_detailed(self, profile: Dict, pair_data: Dict) -> Dict:
        """Detailed token analysis combining profile and pair data"""
        try:

            base_token = pair_data.get('baseToken', {})
            quote_token = pair_data.get('quoteToken', {})
            
            analysis = {
                'url': profile.get('url'),
                'timestamp': datetime.now().isoformat(),
                'token_address': base_token.get('address'),
                'token_name': base_token.get('name'),
                'token_symbol': base_token.get('symbol'),
                'chain': pair_data.get('chainId'),
                'dex_id': pair_data.get('dexId'),
                'pair_address': pair_data.get('pairAddress'),
                'quote_token': quote_token.get('symbol'),
                'creation_time': pair_data.get('pairCreatedAt'),
                'profile_info': {
                    'social_links': profile.get('links', {}),
                    'description': profile.get('description'),
                    'type': profile.get('type'),
                    'tags': profile.get('tags', [])
                },
                'metrics': {},
                'risk_score': 0,
                'flags': [],
                'positive_indicators': []
            }

            # Liquidity Analysis
            liquidity_usd = float(pair_data.get('liquidity', {}).get('usd', 0))
            analysis['metrics']['liquidity_usd'] = liquidity_usd
            
            if liquidity_usd < 10000:
                analysis['risk_score'] += 20
                analysis['flags'].append('Very low liquidity')
            elif liquidity_usd > 100000:
                analysis['risk_score'] -= 10
                analysis['positive_indicators'].append('Strong liquidity')

            # Price Analysis
            price_data = pair_data.get('priceUsd', 0)
            try:
                price_usd = float(price_data)
                analysis['metrics']['price_usd'] = price_usd
                
                # Price changes
                price_change = pair_data.get('priceChange', {})
                analysis['metrics']['price_change_5m'] = float(price_change.get('m5', 0))
                analysis['metrics']['price_change_1h'] = float(price_change.get('h1', 0))
                analysis['metrics']['price_change_24h'] = float(price_change.get('h24', 0))
                
            except (TypeError, ValueError):
                analysis['metrics']['price_usd'] = 0

            # Volume Analysis
            volume_24h = float(pair_data.get('volume', {}).get('h24', 0))
            analysis['metrics']['volume_24h'] = volume_24h

            if volume_24h > 0 and liquidity_usd > 0:
                volume_to_liquidity = volume_24h / liquidity_usd
                analysis['metrics']['volume_to_liquidity_ratio'] = volume_to_liquidity

                if volume_to_liquidity > 5:
                    analysis['flags'].append('Suspicious volume/liquidity ratio')
                    analysis['risk_score'] += 15

            # Profile-based analysis
            websites = [link.get('url') for link in profile.get('links', []) if link.get('type') == 'website']
            if websites:
                analysis['positive_indicators'].append('Has website')
                analysis['risk_score'] -= 5
            
            if profile.get('description'):
                analysis['positive_indicators'].append('Has description')
                analysis['risk_score'] -= 5

            return analysis

        except Exception as e:
            logger.error(f"Error analyzing token: {str(e)}")
            return None

    async def get_token_pairs(self, chain_id: str, token_address: str) -> List[Dict]:
        """Fetch all pairs for a specific token"""
        try:
            endpoint = f"{self.base_url}/tokens/v1/{chain_id}/{token_address}"
            
            async with self.session.get(endpoint) as response:
                if response.status != 200:
                    logger.error(f"Error fetching token pairs: Status {response.status}")
                    return []
                
                pairs = await response.json()
                logger.info(f"Retrieved {len(pairs)} pairs for token {token_address}")
                return pairs

        except Exception as e:
            logger.error(f"Error fetching token pairs: {str(e)}")
            return []

    async def scan_new_tokens(self):
        """Continuously scan for new tokens using profile endpoint"""
        while True:
            try:
                # Get latest token profiles
                profiles = await self.get_latest_token_profiles()
                
                if profiles:
                    logger.info(f"Processing {len(profiles)} token profiles")
                    
                    for profile in profiles:
                        chain_id = profile.get('chainId')
                        token_address = profile.get('tokenAddress')
                        
                        if chain_id and token_address:
                            pairs = await self.get_token_pairs(chain_id, token_address)
                            
                            # Use the first/main pair for analysis
                            if pairs and len(pairs) > 0:
                                pair_data = pairs[0]  # Use first pair as primary
                                pair_address = pair_data.get('pairAddress')
                                
                                if pair_address and pair_address not in self.analyzed_tokens:
                                    analysis = await self.analyze_token_detailed(profile, pair_data)
                                    if analysis:
                                        self.analyzed_tokens.add(pair_address)
                                        
                                        if self.evaluate_opportunity(analysis):
                                            self.promising_tokens.append(analysis)
                                            await self.alert_promising_token(analysis)
                                            await self.export_results()

                            # Respect rate limiting
                            await asyncio.sleep(self.profile_interval)
                
                # Sleep before next scan
                logger.info("Completed scan cycle, waiting before next scan...")
                await asyncio.sleep(60)  # Wait 1 minute before next scan

            except Exception as e:
                logger.error(f"Error during token scanning: {str(e)}")
                await asyncio.sleep(30)

    def evaluate_opportunity(self, analysis: Dict) -> bool:
        """Evaluate if token presents a potential opportunity"""
        if not analysis or analysis['risk_score'] > 50:
            return False

        metrics = analysis.get('metrics', {})
        profile_info = analysis.get('profile_info', {})
        
        required_conditions = [
            metrics.get('liquidity_usd', 0) > 50000,
            metrics.get('volume_24h', 0) > 10000,
            len(analysis['flags']) < 3,
            bool(profile_info.get('social_links', {})),  # Has some social links
            bool(profile_info.get('description')),  # Has description
        ]

        return all(required_conditions)

    async def alert_promising_token(self, analysis: Dict):
        """Alert about promising token opportunities"""
        alert_msg = f"""
🔍 Promising Token Detected!
Dexscreener: {analysis['url']}
Name: {analysis['token_name']} ({analysis['token_symbol']})
Address: {analysis['token_address']}
Chain: {analysis['chain']}
DEX: {analysis['dex_id']}
Pair Address: {analysis['pair_address']}
Price: ${analysis.get('metrics', {}).get('price_usd', 0):,.8f}
Liquidity: ${analysis.get('metrics', {}).get('liquidity_usd', 0):,.2f}
24h Volume: ${analysis.get('metrics', {}).get('volume_24h', 0):,.2f}
Description: {analysis.get('profile_info', {}).get('description', 'N/A')}
Type: {analysis.get('profile_info', {}).get('type', 'N/A')}
Tags: {', '.join(analysis.get('profile_info', {}).get('tags', []))}
Risk Score: {analysis['risk_score']}
Positive Indicators: {', '.join(analysis['positive_indicators'])}
Flags: {', '.join(analysis['flags'])}
        """
        logger.info(alert_msg)

    async def export_results(self):
        """Export analysis results to file"""
        try:
            df = pd.DataFrame(self.promising_tokens)
            df.to_csv('promising_tokens.csv', index=False)
            
            with open('token_analysis.json', 'w') as f:
                json.dump({
                    'last_updated': datetime.now().isoformat(),
                    'total_tokens_analyzed': len(self.analyzed_tokens),
                    'promising_tokens': self.promising_tokens
                }, f, indent=2)
        except Exception as e:
            logger.error(f"Error exporting results: {str(e)}")

async def main():
    scanner = DexTokenScanner()
    await scanner.initialize()
    
    try:
        logger.info("Starting token scanner using profile endpoint...")
        await scanner.scan_new_tokens()
    
    except KeyboardInterrupt:
        logger.info("Scanner stopped by user")
    finally:
        await scanner.close()

if __name__ == "__main__":
    asyncio.run(main())