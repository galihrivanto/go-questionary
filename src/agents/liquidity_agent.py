from .base_agent import BaseAgent

class LiquidityAgent(BaseAgent):
    def get_prompt_template(self) -> str:
        return """
        Analyze the liquidity metrics for the following token:
        {token_data}
        
        Focus on:
        1. Liquidity depth and stability
        2. Volume to liquidity ratio
        3. Liquidity provider concentration
        4. Historical liquidity patterns
        
        Provide a risk assessment and recommendation.
        """
    
    async def analyze(self, token_data: dict) -> dict:
        analysis = await self.chain.arun(token_data=token_data)
        return {
            "agent_type": "liquidity",
            "analysis": analysis,
            "risk_score": self._calculate_risk_score(token_data)
        }
    
    def _calculate_risk_score(self, token_data: dict) -> float:
        # Implement specific liquidity risk scoring logic
        pass